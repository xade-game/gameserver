package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/mattn/go-pubsub"
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
	ps     *pubsub.PubSub
}

func RandomClient(ctx context.Context, server chan []byte) *Client {
	a := rand.Intn(100)
	if a < 40 {
		return NewClient(ctx, server)
	}
	return nil
}

func NewClient(ctx context.Context, server chan []byte) *Client {
	ctx, cancel := context.WithCancel(ctx)
	id := rand.Intn(200)
	conn := NewDummyConn(server)
	c := &Client{
		id:     id,
		status: 0,
		conn:   conn,
		ps:     pubsub.New(),
	}

	c.ps.Sub(func(data []byte) {
		c.status = started
		msec := rand.Intn(50) * 100
		fmt.Printf("client(%d) will die after %d milli second\n", c.id, msec)
		time.Sleep(time.Duration(msec) * time.Millisecond)
		req := &CommandData{
			ClientId: c.id,
			Command:  "update",
			Status:   "dead",
		}
		jsonData, _ := json.Marshal(req)
		c.conn.SendServer(jsonData)
		cancel()
	})
	go c.DataReceive(ctx)
	return c
}

func (c *Client) SendData(data []byte) {
	c.conn.Client() <- data
}

func (c *Client) DataReceive(ctx context.Context) {
	for {
		select {
		case data := <-c.conn.Client():
			c.ps.Pub(data)
		case <-ctx.Done():
			return
		}
	}
}
