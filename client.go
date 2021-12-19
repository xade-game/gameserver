package main

import (
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

type Client struct {
	id     int
	status ClientState
	send   chan []byte
	recv   chan []byte
}

func RandomClient(input chan []byte) *Client {
	a := rand.Intn(100)
	if a < 40 {
		id := rand.Intn(200)
		recv := make(chan []byte, 30)
		c := &Client{
			id:     id,
			status: 0,
			send:   input,
			recv:   recv,
		}

		go c.EventHandler()

		return c
	}
	return nil
}

func (c *Client) EventHandler() {
	for data := range c.recv {
		fmt.Printf("client(%d): recevied: %s\n", c.id, string(data))

		switch string(data) {
		case "opended":
			c.status = started
			c.Run()
		}
	}
}

func (c *Client) Run() {
	msec := rand.Intn(50) * 100
	fmt.Printf("client(%d) will die after %d milli second\n", c.id, msec)
	time.Sleep(time.Duration(msec) * time.Millisecond)
	data := []byte(fmt.Sprintf("{\"id\": %d, \"status\": \"dead\"}", c.id))
	c.send <- data
}

func ClientGenerator() {
}
