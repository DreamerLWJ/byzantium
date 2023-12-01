package process_checker

type EventType int8

const (
	ProcessDown  EventType = 1
	ProcessAlive EventType = 2
	PortDown     EventType = 3
	PortAlive    EventType = 4
)

type ProcessCheckEvent struct {
	EventType EventType
	Pid       int
	Port      int
}

type ProcessDownEventListener interface {
	ProcessDownEvent(ProcessCheckEvent)
}

// ProcessDownEventFunc func called when checking process finished or being killed
type ProcessDownEventFunc func(ProcessCheckEvent)

func (p ProcessDownEventFunc) ProcessDownEvent(event ProcessCheckEvent) {
	p.ProcessDownEvent(event)
}

func SubscribeDownEvent(f ProcessDownEventFunc) Option {
	return func(p *ProcessChecker) {
		if f == nil {
			return
		}
		p.downEventFuncs = append(p.downEventFuncs, f)
	}
}

type ProcessAliveEventListener interface {
	ProcessAliveEvent(ProcessCheckEvent)
}

// ProcessAliveEventFunc func called when checking process alive in check interval
type ProcessAliveEventFunc func(ProcessCheckEvent)

func (p ProcessAliveEventFunc) ProcessAliveEvent(event ProcessCheckEvent) {
	p.ProcessAliveEvent(event)
}

func SubscribeAliveEvent(f ProcessAliveEventFunc) Option {
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
