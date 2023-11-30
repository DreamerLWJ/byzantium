package utils

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/atomic"
)

type AsyncFutureTaskFunc func() (any, error)

type AsyncFutureTask interface {
	FutureTaskFunc() (any, error)
}

type AsyncTaskFunc func()

type AsyncTask interface {
	AsyncTaskFunc()
}

type AsyncControlPlane struct {
	wg          sync.WaitGroup
	tasks       map[uint32]AsyncTaskFunc
	futureTasks map[uint32]AsyncFutureTaskFunc

	futureTaskChs map[uint32]chan<- futureChElem
	isCalled      *atomic.Bool
	isRunning     *atomic.Bool
	isTerminated  *atomic.Bool
	taskIDGen     *atomic.Uint32
}

func NewAsyncControlPlane() *AsyncControlPlane {
	return &AsyncControlPlane{
		tasks:         make(map[uint32]AsyncTaskFunc),
		futureTasks:   make(map[uint32]AsyncFutureTaskFunc),
		futureTaskChs: make(map[uint32]chan<- futureChElem),
		isCalled:      atomic.NewBool(false),
		isRunning:     atomic.NewBool(false),
		isTerminated:  atomic.NewBool(false),
		taskIDGen:     atomic.NewUint32(0),
	}
}

func (a *AsyncControlPlane) genTaskID() uint32 {
	return a.taskIDGen.Inc()
}

func (a *AsyncControlPlane) Task(task AsyncTask) *AsyncControlPlane {
	a.TaskFn(task.AsyncTaskFunc)
	return a
}

func (a *AsyncControlPlane) TaskFn(fn AsyncTaskFunc) *AsyncControlPlane {
	taskID := a.genTaskID()
	a.tasks[taskID] = fn
	return a
}

func (a *AsyncControlPlane) FutureTask(task AsyncFutureTask) Future {
	return a.FutureTaskFn(task.FutureTaskFunc)
}

func (a *AsyncControlPlane) FutureTaskFn(fn AsyncFutureTaskFunc) Future {
	taskID := a.genTaskID()
	a.futureTasks[taskID] = fn
	ch := make(chan futureChElem, 1)
	a.futureTaskChs[taskID] = ch
	return newFuture(ch, a)
}

// Sync start to run all binding task, please don't repeated call
func (a *AsyncControlPlane) Sync(ctx context.Context) error {
	// not allow repeated call
	if isCalled := a.isCalled.Swap(true); isCalled {
		return errors.Errorf("not allow repeated call")
	}

	// normally, if caller always once call, it won't return error below
	if a.isRunning.Load() {
		return errors.Errorf("already running")
	}
	if len(a.tasks) <= 0 && len(a.futureTasks) <= 0 {
		return nil
	}

	a.isRunning.Store(true)
	for _, fn := range a.tasks {
		taskFn := fn
		GoWithWg(ctx, func() {
			a.runTask(taskFn)
		}, &a.wg)
	}

	for id, fn := range a.futureTasks {
		taskFn := fn
		taskID := id
		GoWithWg(ctx, func() {
			a.runFutureTask(taskID, taskFn)
		}, &a.wg)
	}

	a.wg.Wait()

	a.isRunning.Store(false)
	a.isTerminated.Store(true)
	return nil
}

func (a *AsyncControlPlane) runTask(fn AsyncTaskFunc) {
	fn()
}

func (a *AsyncControlPlane) runFutureTask(taskID uint32, fn AsyncFutureTaskFunc) {
	res, err := fn()
	a.futureTaskChs[taskID] <- futureChElem{
		res: res,
		err: err,
	}
	close(a.futureTaskChs[taskID])
}

type futureChElem struct {
	res any
	err error
}

type Future struct {
	p  *AsyncControlPlane
	ch <-chan futureChElem
}

func newFuture(ch chan futureChElem, p *AsyncControlPlane) Future {
	return Future{ch: ch, p: p}
}

func (f *Future) Get() (res any, err error) {
	// avoid deadlock calling
	if !f.p.isRunning.Load() && !f.p.isTerminated.Load() {
		return nil, errors.Errorf("plane not running or terminated")
	}

	elem := <-f.ch
	return elem.res, elem.err
}
