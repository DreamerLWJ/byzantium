package process

// TcpConnStatus tcp connection status
type TcpConnStatus string

const (
	Listen      TcpConnStatus = "LISTEN"
	ESTABLISHED TcpConnStatus = "ESTABLISHED"
	CloseWait   TcpConnStatus = "CLOSE_WAIT"
	TimeWait    TcpConnStatus = "TIME_WAIT"
	Closing     TcpConnStatus = "CLOSING"
	SynSent     TcpConnStatus = "SYN_SENT"
	SynRecv     TcpConnStatus = "SYN_RECV"
	FinWait1    TcpConnStatus = "FIN_WAIT1"
	FinWait2    TcpConnStatus = "FIN_WAIT2"
	LastAck     TcpConnStatus = "LAST_ACK"
)
