package process_checker

type EventType int8

const (
	ProcessDown    EventType = 1
	ProcessAlive   EventType = 2
	PortDown       EventType = 3
	PortAlive      EventType = 4
	PortPidChanged EventType = 5
)

type ProcessCheckEvent struct {
	EventType EventType
	OldPid    int
	Pid       int
	Port      int
}
