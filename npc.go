package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/myoan/snake/api"
)

type NonPlayerClient struct {
	uuid   string
	stream chan []byte
}

func NewNonPlayerClient() *NonPlayerClient {
	stream := make(chan []byte)

	return &NonPlayerClient{
		uuid:   uuid.NewString(),
		stream: stream,
	}
}

func (c *NonPlayerClient) ID() string {
	return c.uuid
}

func (c *NonPlayerClient) Send(data []byte) error {
	log.Printf("Send to client %s: %s", c.ID(), data)

	var resp api.EventResponse
	err := json.Unmarshal(data, &resp)
	if err != nil {
		log.Printf("[Error] decode(%s): %v", c.ID(), err)
		return err
	}

	for _, p := range resp.Body.Players {
		if p.ID == c.ID() {
			curX := p.X
			curY := p.Y

			curX++
			curY++

			req := EventRequest{
				UUID:      c.ID(),
				Eventtype: 0,
				X:         curX,
				Y:         curY,
				Theta:     -20,
			}

			data, err := json.Marshal(req)
			if err != nil {
				log.Printf("[Error] write(%s): %v", c.ID(), err)
				return err
			}

			c.stream <- data
			return nil
		}
	}

	return nil
}

func (c *NonPlayerClient) Close() {
	close(c.stream)
}

func (c *NonPlayerClient) Stream() chan []byte {
	return c.stream
}
