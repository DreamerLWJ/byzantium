package process_checker

import (
	"log"
	"time"
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

func Listener(l ProcessCheckEventListener) Option {
	return func(p *ProcessChecker) {
		if l == nil {
			return
		}
		p.aliveEventFuncs = append(p.aliveEventFuncs, l.ProcessAliveEvent)
		p.downEventFuncs = append(p.downEventFuncs, l.ProcessAliveEvent)
	}
}
