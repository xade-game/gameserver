package main

import (
	"encoding/json"

	"github.com/myoan/snake/api"
)

type Player struct {
	x      int
	y      int
	done   chan struct{}
	Client Client
}

func (p *Player) ID() string {
	return p.Client.ID()
}
func NewPlayer(client Client, stream <-chan []byte, x, y int) *Player {
	p := &Player{
		x:      x,
		y:      y,
		Client: client,
	}
	go p.run(stream)
	return p
}

func (p *Player) Finish() {
	p.Client.Close()
}

func (p *Player) Send(status int, players []*Player) error {
	playersProtocol := make([]api.PlayerResponse, len(players))
	for i, player := range players {
		playersProtocol[i] = api.PlayerResponse{
			ID: player.ID(),
			X:  player.x,
			Y:  player.y,
		}
	}

	resp := &api.EventResponse{
		Status: status,
		Body: api.ResponseBody{
			Board:   []int{},
			Width:   0,
			Height:  0,
			Players: playersProtocol,
		},
	}

	bytes, _ := json.Marshal(&resp)
	return p.Client.Send(bytes)
}

func (p *Player) Move() error {
	return nil
}

func (p *Player) ChangeDirection(direction int) {
}

func (p *Player) run(stream <-chan []byte) {
	for {
		select {
		case <-p.done:
			return
		case msg := <-stream:
			var req api.EventRequest
			json.Unmarshal(msg, &req)

			p.ChangeDirection(req.Key)
		}
	}
}
