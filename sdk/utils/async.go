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
	isRunning     *atomic.Bool
	taskIDGen     *atomic.Uint32
}

func NewAsyncControlPlane() *AsyncControlPlane {
	return &AsyncControlPlane{
		isRunning: atomic.NewBool(false),
		taskIDGen: atomic.NewUint32(0),
	}
}

func (a *AsyncControlPlane) genTaskID() uint32 {
	return a.taskIDGen.Inc()
}

func (a *AsyncControlPlane) Task(task AsyncTask) {
	a.TaskFn(task.AsyncTaskFunc)
}

func (a *AsyncControlPlane) TaskFn(fn AsyncTaskFunc) {
	taskID := a.genTaskID()
	a.tasks[taskID] = fn
}

func (a *AsyncControlPlane) FutureTask(task AsyncFutureTask) Future {
	return a.FutureTaskFn(task.FutureTaskFunc)
}

func (a *AsyncControlPlane) FutureTaskFn(fn AsyncFutureTaskFunc) Future {
	taskID := a.genTaskID()
	a.futureTasks[taskID] = fn
	ch := make(chan futureChElem, 1)
	a.futureTaskChs[taskID] = ch
	return newFuture(ch)
}

func (a *AsyncControlPlane) Sync(ctx context.Context) error {
	if a.isRunning.Load() {
		return errors.Errorf("already running")
	}
	if len(a.tasks) <= 0 {
		return nil
	}
	for _, fn := range a.tasks {
		taskFn := fn
		GoWithWg(ctx, func() {
			a.runTask(taskFn)
		}, &a.wg)
	}

	for taskID, fn := range a.futureTasks {
		taskFn := fn
		GoWithWg(ctx, func() {
			a.runFutureTask(taskID, taskFn)
		}, &a.wg)
	}

	a.wg.Wait()
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
	ch <-chan futureChElem
}

func newFuture(ch chan futureChElem) Future {
	return Future{ch: ch}
}

func (f *Future) Get() (res any, err error) {
	elem := <-f.ch
	return elem.res, elem.err
}
