package system

type Client interface {
	ID() string
	Send(data []byte) error
	Close()
	Stream() chan []byte
}

type TriggerArgument struct {
	EventType int
	Client    Client
}

const (
	EventClientConnect = iota
	EventClientFinish
	EventClientRestart
)
