package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type ClientState int

const (
	_ = iota
	registered
	started
	dead
)

type IConn interface {
	SendServer(data []byte)
	Client() chan []byte
}

type DummyConn struct {
	server chan []byte
	client chan []byte
}

func NewDummyConn(server chan []byte) *DummyConn {
	client := make(chan []byte, 30)
	conn := &DummyConn{
		server: server,
		client: client,
	}

	return conn
}

func (conn *DummyConn) SendServer(data []byte) {
	conn.server <- data
}

func (conn *DummyConn) Client() chan []byte {
	return conn.client
}

type Client struct {
	id     int
	status ClientState
	conn   IConn
}

func RandomClient(server chan []byte) *Client {
	a := rand.Intn(100)
	if a < 40 {
		return NewClient(server)
	}
	return nil
}

func NewClient(server chan []byte) *Client {
	id := rand.Intn(200)
	conn := NewDummyConn(server)
	c := &Client{
		id:     id,
		status: 0,
		conn:   conn,
	}

	go c.EventHandler()

	return c
}

func (c *Client) SendData(data []byte) {
	c.conn.Client() <- data
}

func (c *Client) EventHandler() {
	for data := range c.conn.Client() {
		fmt.Printf("client(%d): recevied: %s\n", c.id, string(data))

		game := &GameStateData{}
		json.Unmarshal(data, game)

		switch game.Status {
		case "opened":
			c.status = started
			msec := rand.Intn(50) * 100
			fmt.Printf("client(%d) will die after %d milli second\n", c.id, msec)
			time.Sleep(time.Duration(msec) * time.Millisecond)
			data := &CommandData{
				ClientId: c.id,
				Command:  "update",
				Status:   "dead",
			}
			jsonData, _ := json.Marshal(data)
			c.conn.SendServer(jsonData)
		}
	}
}
