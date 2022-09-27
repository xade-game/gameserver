package main

import (
	"context"
	"encoding/json"

	"github.com/xade-game/gameserver/api"
	"github.com/xade-game/gameserver/system"
)

const (
	PlayerAlive = iota
	PlayerDead
)

type Player struct {
	x      int
	y      int
	theta  int
	done   chan struct{}
	Client system.Client
	Status int
}

var ctx context.Context

func NewPlayer(client system.Client, stream <-chan []byte, x, y int) *Player {
	p := &Player{
		x:      x,
		y:      y,
		Client: client,
		Status: PlayerAlive,
	}
	go p.run(stream)
	return p
}

func (p *Player) ID() string {
	return p.Client.ID()
}

func (p *Player) Finish() {
	p.Client.Close()
}

func (p *Player) Send(status int, board *Board, players []*Player) error {
	playersProtocol := make([]api.PlayerResponse, len(players))
	for i, player := range players {
		playersProtocol[i] = api.PlayerResponse{
			ID:        player.ID(),
			X:         player.x,
			Y:         player.y,
			Direction: player.theta,
		}
	}

	resp := &api.EventResponse{
		Status: status,
		Body: api.ResponseBody{
			Board:   board.ToArray(),
			Width:   GameCellWidth,
			Height:  GameCellHeight,
			Players: playersProtocol,
		},
	}

	bytes, _ := json.Marshal(&resp)
	return p.Client.Send(bytes)
}

func (p *Player) Move(x, y, theta int) {
	p.x = x
	p.y = y
	p.theta = theta
}

func (p *Player) run(stream <-chan []byte) {
	for {
		select {
		case <-p.done:
			return
		case msg := <-stream:
			var req api.EventRequest
			json.Unmarshal(msg, &req)

			p.Move(req.X, req.Y, req.Theta)
		}
	}
}

func FindByID(id string) (*Player, error) {
	return nil, nil
}
