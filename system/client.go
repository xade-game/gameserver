package system

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

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

type WebClient struct {
	Uuid      string
	stream    chan []byte
	conn      *websocket.Conn
	observers []Observer
	mu        sync.Mutex
}

func NewWebClient(c *websocket.Conn) *WebClient {
	stream := make(chan []byte)
	obs := make([]Observer, 0)
	return &WebClient{
		Uuid:      uuid.NewString(),
		stream:    stream,
		conn:      c,
		observers: obs,
	}
}

func (c *WebClient) AddObserver(o Observer) {
	c.observers = append(c.observers, o)
}

func (c *WebClient) Notify(tp int) {
	for _, o := range c.observers {
		data := TriggerArgument{
			EventType: tp,
			Client:    c,
		}
		o.Update(data)
	}
}

func (c *WebClient) ID() string {
	return c.Uuid
}

func (c *WebClient) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	log.Printf("Send to client %s: %s", c.ID(), data)
	err := c.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("[Error] write(%s): %v", c.ID(), err)
		return err
	}
	return nil
}

func (c *WebClient) Stream() chan []byte {
	return c.stream
}

func (c *WebClient) Run() {
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("[Error] read: ", err)
			log.Printf("message type: %d", mt)
			close(c.Stream())
			c.Close()
			return
		}
		log.Printf("recv: %s", message)

		c.Stream() <- message
	}
}

func (c *WebClient) Close() {
	log.Printf("Close client %s", c.ID())
	c.Notify(EventClientFinish)
	c.conn.Close()
}
