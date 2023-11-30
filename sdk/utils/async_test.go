package utils

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTask struct {
	Start int
}

func (t testTask) AsyncTaskFunc() {
	for i := t.Start; i < t.Start+10; i++ {
		time.Sleep(time.Second)
		fmt.Println("testTask", i)
	}
}

type testFutureTask struct {
	Start int
}

func (t testFutureTask) FutureTaskFunc() (any, error) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second + time.Second*time.Duration(rand.Intn(1)))
		t.Start += 1
		fmt.Println("testFutureTask", t.Start)
	}
	return t.Start, nil
}

func TestAsyncControlPlane_Task(t *testing.T) {
	plane := NewAsyncControlPlane()

	t1 := testTask{Start: 10}
	t2 := testTask{Start: 1000}
	t3 := testTask{Start: 1000}

	plane.Task(t1).Task(t2).Task(t3)

	err := plane.Sync(context.Background())
	assert.Nil(t, err)
}

func TestAsyncControlPlane_TaskFn(t *testing.T) {
	plane := NewAsyncControlPlane()

	t1 := testTask{Start: 10}
	t2 := testTask{Start: 100}
	t3 := testTask{Start: 1000}

	plane.TaskFn(t1.AsyncTaskFunc).TaskFn(t2.AsyncTaskFunc).TaskFn(t3.AsyncTaskFunc)

	err := plane.Sync(context.Background())
	assert.Nil(t, err)
}

func TestAsyncControlPlane_FutureTask(t *testing.T) {
	plane := NewAsyncControlPlane()
	t1 := testFutureTask{Start: 10}
	t2 := testFutureTask{Start: 100}
	t3 := testFutureTask{Start: 1000}

	f1 := plane.FutureTask(t1)
	f2 := plane.FutureTask(t2)
	f3 := plane.FutureTask(t3)

	err := plane.Sync(context.Background())
	assert.Nil(t, err)
	fmt.Println(f1.Get())
	fmt.Println(f2.Get())
	fmt.Println(f3.Get())
}

func TestAsyncControlPlane_FutureTaskFn(t *testing.T) {
	plane := NewAsyncControlPlane()
	t1 := testFutureTask{Start: 10}
	t2 := testFutureTask{Start: 100}
	t3 := testFutureTask{Start: 1000}

	f1 := plane.FutureTaskFn(t1.FutureTaskFunc)
	f2 := plane.FutureTaskFn(t2.FutureTaskFunc)
	f3 := plane.FutureTaskFn(t3.FutureTaskFunc)

	err := plane.Sync(context.Background())
	assert.Nil(t, err)
	fmt.Println(f1.Get())
	fmt.Println(f2.Get())
	fmt.Println(f3.Get())
}

func TestAsyncControlPlane_Sync(t *testing.T) {
	plane := NewAsyncControlPlane()
	t1 := testTask{Start: 99}
	t2 := testFutureTask{Start: 999}
	t3 := testTask{Start: 9999}
	t4 := testFutureTask{Start: 99999}
	plane.Task(t1)
	f2 := plane.FutureTask(t2)
	plane.Task(t3)
	f4 := plane.FutureTask(t4)

	err := plane.Sync(context.Background())
	assert.Nil(t, err)
	res, err := f2.Get()
	assert.Nil(t, err)
	assert.NotNil(t, res)
	t.Logf("f2.Get (res:%+v,err:%s)", res, err)
	res, err = f4.Get()
	assert.Nil(t, err)
	assert.NotNil(t, res)
	t.Logf("f2.Get (res:%+v,err:%s)", res, err)
}

func TestAsyncControlPlane_Deadlock(t *testing.T) {
	plane := NewAsyncControlPlane()
	t1 := testFutureTask{Start: 99}
	f1 := plane.FutureTask(t1)

	res, err := f1.Get()
	t.Logf("f1.Get (res:%+v,err:%s)", res, err)
	assert.Nil(t, res)
	assert.NotNil(t, err)

	err = plane.Sync(context.Background())
	assert.Nil(t, err)

	res, err = f1.Get()
	t.Logf("f1.Get (res:%+v,err:%s)", res, err)
	assert.NotNil(t, res)
	assert.Nil(t, err)
}

func TestAsyncControlPlane_RepeatedCall(t *testing.T) {
	plane := NewAsyncControlPlane()
	t1 := testTask{Start: 99}
	plane.Task(t1)

	err := plane.Sync(context.Background())
	assert.Nil(t, err)

	err = plane.Sync(context.Background())
	assert.NotNil(t, err)
}
