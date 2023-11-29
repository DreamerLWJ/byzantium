package availability

import (
	"byzantium-availability/sdk/process"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type CheckerType int8

const (
	PortChecker CheckerType = 1
	PidChecker  CheckerType = 2
)

type Option func(p *ProcessChecker)

func SetLogger(logger *log.Logger) Option {
	return func(p *ProcessChecker) {

	}
}

// CheckInterval specific checker interval
func CheckInterval(interval time.Duration) Option {
	return func(p *ProcessChecker) {
		p.checkInterval = interval
	}
}

type ProcessDownEventListener interface {
	ProcessDownEvent(pid, port int)
}

// ProcessDownEventFunc func called when checking process finished or being killed
type ProcessDownEventFunc func(pid, port int)

func (p ProcessDownEventFunc) ProcessDownEvent(pid, port int) {
	p.ProcessDownEvent(pid, port)
}

func ListenDownEvent(f ProcessDownEventFunc) Option {
	return func(p *ProcessChecker) {
		if f == nil {
			return
		}
		p.downEventFuncs = append(p.downEventFuncs, f)
	}
}

type ProcessAliveEventListener interface {
	ProcessAliveEvent(pid, port int)
}

// ProcessAliveEventFunc func called when checking process alive in check interval
type ProcessAliveEventFunc func(pid, port int)

func (p ProcessAliveEventFunc) ProcessAliveEvent(pid, port int) {
	p.ProcessAliveEvent(pid, port)
}

func ListenAliveEvent(f ProcessAliveEventFunc) Option {
	return func(p *ProcessChecker) {
		if f == nil {
			return
		}
		p.aliveEventFuncs = append(p.aliveEventFuncs, f)
	}
}

type ProcessCheckEventListener interface {
	ProcessAliveEventListener
	ProcessDownEventListener
}

func Listener(l ProcessCheckEventListener) Option {
	return func(p *ProcessChecker) {
		if l == nil {
			return
		}
		p.aliveEventFuncs = append(p.aliveEventFuncs, l.ProcessAliveEvent)
		p.downEventFuncs = append(p.downEventFuncs, l.ProcessAliveEvent)
	}
}

// ProcessChecker check if process exist
type ProcessChecker struct {
	pid           int           // specific pid to check
	port          int           // check port pid alive if pid not specified
	checkerType   CheckerType   // checker type
	checkInterval time.Duration // check process interval

	downEventFuncs  []ProcessDownEventFunc
	aliveEventFuncs []ProcessAliveEventFunc

	mu           sync.RWMutex
	isRunning    bool
	isTerminated bool
}

func (p *ProcessChecker) initDefaultValue() {
	p.checkInterval = time.Second * 5
}

func NewPortProcessChecker(port int, opts ...Option) (*ProcessChecker, error) {
	if port <= 0 {
		return &ProcessChecker{}, errors.Errorf("port <= 0")
	}

	p := &ProcessChecker{port: port, checkerType: PortChecker}
	p.initDefaultValue()
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
}

func NewPidProcessChecker(pid int, opts ...Option) (*ProcessChecker, error) {
	if pid <= 0 {
		return &ProcessChecker{}, errors.Errorf("pid <= 0")
	}
	p := &ProcessChecker{pid: pid, checkerType: PidChecker}
	p.initDefaultValue()
	for _, opt := range opts {
		opt(p)
	}
	return p, nil
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

func (p *ProcessChecker) Run() error {
	err := p.preCheck()
	if err != nil {
		return err
	}

	p.doCheck()
	return nil
}

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

func (p *ProcessChecker) doCheck() {
	if p.pid > 0 {
		alive, err := process.IsProcessAlive(p.pid)
		if err != nil {
			// TODO
		}

		wg := sync.WaitGroup{}
		if alive {
			wg.Add(len(p.aliveEventFuncs))
			for _, f := range p.aliveEventFuncs {
				go func(fn ProcessAliveEventFunc) {
					fn(p.pid, p.port)
				}(f)
			}
		} else {
			wg.Add(len(p.downEventFuncs))
			for _, f := range p.downEventFuncs {
				go func(fn ProcessDownEventFunc) {
					fn(p.pid, p.port)
				}(f)
			}
		}
		time.Sleep(p.checkInterval)
		wg.Wait()
	}
}
