package process_checker

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEventListener struct {
}

func (t testEventListener) ProcessAliveEvent(ctx context.Context, event ProcessCheckEvent) {
	logHead := "ProcessAliveEvent|"
	fmt.Printf(logHead+"event:%+v\n", event)
}

func (t testEventListener) ProcessDownEvent(ctx context.Context, event ProcessCheckEvent) {
	logHead := "ProcessDownEvent|"
	fmt.Printf(logHead+"event:%+v\n", event)
}

func (t testEventListener) PortPidChangedEvent(ctx context.Context, event ProcessCheckEvent) {
	logHead := "PortPidChangedEvent|"
	fmt.Printf(logHead+"event:%+v\n", event)
}

func TestPortProcessChecker(t *testing.T) {
	checker, err := NewPortProcessChecker(8080, ListenEvent(testEventListener{}))
	assert.Nil(t, err)
	err = checker.Run()
	if err != nil {
		return
	}
}
