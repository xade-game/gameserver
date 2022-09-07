package cambrian

import (
	"time"

	"github.com/gorilla/websocket"
)

type Cambrian struct {
	onConnect    func(*WebSocketClient, interface{})
	onDisconnect func(*WebSocketClient, interface{})
	onMessage    func(Request, interface{})
}

func New() *Cambrian {
	return &Cambrian{}
}

func (c *Cambrian) Start(uri string) {
}

func (c *Cambrian) RegisterPeriodic(d time.Duration, fn func(Request)) {
	tick := time.NewTicker(d)
	go func() {
		for range tick.C {
			fn(&WebSocketMessage{})
		}
	}()
}

func (c *Cambrian) RegisterWebsocketConnect(fn func(*WebSocketClient, interface{})) {
	c.onConnect = fn
}

func (c *Cambrian) RegisterWebsocketDisconnect(fn func(*WebSocketClient, interface{})) {
	c.onDisconnect = fn
}

func (c *Cambrian) RegisterWebsocketMessage(fn func(Request, interface{})) {
	c.onMessage = fn
}

func (c *Cambrian) AddWebsocketClient(conn *websocket.Conn, obj interface{}) {
	client := newWebSocketClient(conn)
	c.onConnect(client, obj)

	go func() {
		for {
			select {
			case msg := <-client.stream:
				c.onMessage(&WebSocketMessage{
					id:   client.ID(),
					Conn: client.conn,
					body: msg,
				}, obj)
			case <-client.done:
				c.onDisconnect(client, obj)
			}
		}
	}()
}

func NewPeriodicRunner() chan struct{} {
	tick := time.NewTicker(time.Second)
	stream := make(chan struct{})
	go func() {
		for range tick.C {
			stream <- struct{}{}
		}
	}()
	return stream
}

type Request interface {
	ID() string
	Body() []byte
}

type WebSocketMessage struct {
	id   string
	Conn *websocket.Conn
	body []byte
}

func (msg *WebSocketMessage) ID() string {
	return msg.id
}

func (msg *WebSocketMessage) Body() []byte {
	return msg.body
}

func NewWebSockRunner(conn *websocket.Conn) chan []byte {
	client := newWebSocketClient(conn)
	return client.stream
}

func NewCustomRunner(stream chan struct{}) chan struct{} {
	return stream
}
