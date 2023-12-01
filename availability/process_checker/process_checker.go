package process_checker

import (
	"context"
	"sync"
	"time"

	"github.com/DreamerLWJ/byzantium/sdk/process"
	"github.com/DreamerLWJ/byzantium/sdk/utils"
	"github.com/pkg/errors"
)

const (
	_defaultCheckInterval = time.Second * 5
)

type CheckerType int8

const (
	PortChecker CheckerType = 1
	PidChecker  CheckerType = 2
)

// ProcessChecker check if process exist
type ProcessChecker struct {
	pid           int           // specific pid to check, or last listen port pid
	port          int           // check port pid alive if pid not specified
	checkerType   CheckerType   // checker type
	checkInterval time.Duration // check process interval

	downEventFuncs      []ProcessDownEventFunc
	aliveEventFuncs     []ProcessAliveEventFunc
	portPidChangedFuncs []PortPidChangedEventFunc

	mu           sync.RWMutex
	isRunning    bool
	isTerminated bool
}

func NewPortProcessChecker(port int, opts ...Option) (*ProcessChecker, error) {
	if port <= 0 {
		return &ProcessChecker{}, errors.Errorf("port <= 0")
	}

	p := &ProcessChecker{
		port:          port,
		checkerType:   PortChecker,
		checkInterval: _defaultCheckInterval,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func NewPidProcessChecker(pid int, opts ...Option) (*ProcessChecker, error) {
	if pid <= 0 {
		return &ProcessChecker{}, errors.Errorf("pid <= 0")
	}
	p := &ProcessChecker{
		pid:           pid,
		checkerType:   PidChecker,
		checkInterval: _defaultCheckInterval,
	}
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

// Run run checker in blocking way
func (p *ProcessChecker) Run() error {
	err := p.preCheck()
	if err != nil {
		return err
	}
	// listen kill event, TODO

	p.mu.Lock()
	p.isRunning = true
	p.mu.Unlock()

	for {
		p.mu.RLock()
		if p.isTerminated {
			p.mu.RUnlock()
			return ErrTerminated
		}
		p.mu.RUnlock()

		p.doCheck()
		time.Sleep(p.checkInterval)
	}
}

// Start run checker in non-blocking way
func (p *ProcessChecker) Start() error {
	err := p.preCheck()
	if err != nil {
		return err
	}
	go func() {
		_ = p.Run()
	}()
	return nil
}

// Terminate terminate checker
func (p *ProcessChecker) Terminate() {
	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
	}()

	p.isTerminated = true
	p.isRunning = false
}

func (p *ProcessChecker) preCheck() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.isTerminated {
		return errors.Errorf("checker already terminated")
	}
	if p.isRunning {
		return errors.Errorf("checker already running")
	}
	return nil
}

func (p *ProcessChecker) doCheck() {
	switch p.checkerType {
	case PortChecker:
		if p.pid != 0 {
			p.checkPidAlive()
		} else {
			p.checkPortListening()
		}
	case PidChecker:
		p.checkPidAlive()
	}
}

func (p *ProcessChecker) checkPidAlive() {
	ctx := context.Background()
	if p.pid <= 0 {
		return
	}
	alive, err := process.IsProcessAlive(p.pid)
	if err != nil {
		// TODO
	}

	if alive {
		p.dispatchEvents(ctx, ProcessCheckEvent{
			EventType: ProcessAlive,
			Pid:       p.pid,
			Port:      p.port,
		})
	} else {
		p.dispatchEvents(ctx, ProcessCheckEvent{
			EventType: ProcessDown,
			Pid:       p.pid,
			Port:      p.port,
		})
	}
}

func (p *ProcessChecker) checkPortListening() {
	ctx := context.Background()
	if p.port <= 0 {
		return
	}

	pids, err := process.GetPortPID(p.port, process.Listen)
	if err != nil {
		// TODO
	}
	if len(pids) <= 0 { // process binding
		p.dispatchEvents(ctx, ProcessCheckEvent{
			EventType: PortDown,
			Pid:       p.pid,
			Port:      p.port,
		})
	} else { // port isn't bound
		if len(pids) != 1 {
			// TODO
		}
		p.pid = pids[0]

		p.dispatchEvents(ctx, ProcessCheckEvent{
			EventType: ProcessAlive,
			Pid:       p.pid,
			Port:      p.port,
		})
	}
}

func (p *ProcessChecker) dispatchEvents(ctx context.Context, event ProcessCheckEvent) {
	ap := utils.NewAsyncControlPlane()
	defer func() {
		err := ap.Sync(ctx)
		if err != nil {
			// TODO
		}
	}()

	switch event.EventType {
	case PortDown, ProcessDown:
		for _, f := range p.downEventFuncs {
			ap.TaskFn(func() {
				f(ctx, event)
			})
		}
	case PortAlive, ProcessAlive:
		for _, f := range p.aliveEventFuncs {
			ap.TaskFn(func() {
				f(ctx, event)
			})
		}
	}
}
