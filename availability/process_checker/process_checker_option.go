package process_checker

import (
	"log"
	"time"
)

type Option func(p *ProcessChecker)

func SetLogger(logger *log.Logger) Option {
	return func(p *ProcessChecker) {
		// TODO
	}
}

// CheckInterval specific checker interval
func CheckInterval(interval time.Duration) Option {
	return func(p *ProcessChecker) {
		p.checkInterval = interval
	}
}

func ListenDownEvent(l ...ProcessDownEventListener) Option {
	return func(p *ProcessChecker) {
		if len(l) <= 0 {
			return
		}
		funcs := make([]ProcessDownEventFunc, 0, len(l))
		for _, listener := range l {
			funcs = append(funcs, listener.ProcessDownEvent)
		}
		p.downEventFuncs = append(p.downEventFuncs, funcs...)
	}
}

func ConsumeDownEvent(f ...ProcessDownEventFunc) Option {
	return func(p *ProcessChecker) {
		if len(f) <= 0 {
			return
		}
		p.downEventFuncs = append(p.downEventFuncs, f...)
	}
}

func ListenAliveEvent(l ...ProcessAliveEventListener) Option {
	return func(p *ProcessChecker) {
		if len(l) <= 0 {
			return
		}
		funcs := make([]ProcessAliveEventFunc, 0, len(l))
		for _, listener := range l {
			funcs = append(funcs, listener.ProcessAliveEvent)
		}
		p.aliveEventFuncs = append(p.aliveEventFuncs, funcs...)
	}
}

func ConsumeAliveEvent(f ...ProcessAliveEventFunc) Option {
	return func(p *ProcessChecker) {
		if len(f) <= 0 {
			return
		}
		p.aliveEventFuncs = append(p.aliveEventFuncs, f...)
	}
}

func ListenEvent(l ProcessCheckEventListener) Option {
	return func(p *ProcessChecker) {
		if l == nil {
			return
		}
		p.aliveEventFuncs = append(p.aliveEventFuncs, l.ProcessAliveEvent)
		p.downEventFuncs = append(p.downEventFuncs, l.ProcessDownEvent)
		p.portPidChangedFuncs = append(p.portPidChangedFuncs, l.PortPidChangedEvent)
	}
}
