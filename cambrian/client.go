package cambrian

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketClientStatus int

const (
	Init WebSocketClientStatus = iota
	Connect
	Closed
)

type WebSocketClient struct {
	status WebSocketClientStatus
	sessID string
	stream chan []byte
	done   chan struct{}
	conn   *websocket.Conn
	mu     sync.Mutex
}

func (c *WebSocketClient) ID() string {
	return c.sessID
}

func (c *WebSocketClient) Send(data []byte) error {
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

func (c *WebSocketClient) Status() WebSocketClientStatus {
	return c.status
}

func (c *WebSocketClient) Stream() chan []byte {
	return c.stream
}

func (c *WebSocketClient) Close() {
	c.status = Closed
	c.conn.Close()
}

func newWebSocketClient(c *websocket.Conn) *WebSocketClient {
	sessID := uuid.NewString()
	stream := make(chan []byte)
	done := make(chan struct{})
	client := &WebSocketClient{
		status: Connect,
		sessID: sessID,
		stream: stream,
		done:   done,
		conn:   c,
	}
	go client.run()
	return client
}

func (c *WebSocketClient) run() {
	for {
		mt, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("readMessage err: %v\n", err)
			c.done <- struct{}{}
			return
		}
		log.Printf("type: %d, msg: %s\n", mt, msg)
		c.stream <- msg
	}
}
