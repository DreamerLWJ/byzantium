package process_checker

import "context"

type ProcessDownEventListener interface {
	ProcessDownEvent(ctx context.Context, event ProcessCheckEvent)
}

// ProcessDownEventFunc func called when checking process finished or being killed
type ProcessDownEventFunc func(ctx context.Context, event ProcessCheckEvent)

func (p ProcessDownEventFunc) ProcessDownEvent(ctx context.Context, event ProcessCheckEvent) {
	p(ctx, event)
}

type ProcessAliveEventListener interface {
	ProcessAliveEvent(ctx context.Context, event ProcessCheckEvent)
}

// ProcessAliveEventFunc func called when checking process alive in check interval
type ProcessAliveEventFunc func(ctx context.Context, event ProcessCheckEvent)

func (p ProcessAliveEventFunc) ProcessAliveEvent(ctx context.Context, event ProcessCheckEvent) {
	p(ctx, event)
}

type PortPidChangedEventListener interface {
	PortPidChangedEvent(ctx context.Context, event ProcessCheckEvent)
}

type PortPidChangedEventFunc func(ctx context.Context, event ProcessCheckEvent)

func (p PortPidChangedEventFunc) PortPidChangedEvent(ctx context.Context, event ProcessCheckEvent) {
	p(ctx, event)
}

type ProcessCheckEventListener interface {
	ProcessAliveEventListener
	ProcessDownEventListener
	PortPidChangedEventListener
}
