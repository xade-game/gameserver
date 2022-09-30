package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

	"github.com/xade-game/gameserver/api"
	"github.com/xade-game/gameserver/system"
)

const (
	PlayerAlive = iota
	PlayerDead
)

type Player struct {
	x         int
	y         int
	theta     int
	Client    system.Client
	Status    int
	direction int
	size      int
}

func NewPlayer(client system.Client, stream <-chan []byte, w, h int) *Player {
	x := rand.Intn(w)
	y := rand.Intn(h)
	d := rand.Intn(4)

	p := &Player{
		x:         x,
		y:         y,
		direction: d,
		size:      3,
		Client:    client,
		Status:    PlayerAlive,
	}
	return p
}

func (p *Player) ID() string {
	return p.Client.ID()
}

func (p *Player) Finish() {
	p.Status = 1
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

func (p *Player) GenerateSnake(board *Board) {
	log.Printf("GenerateSnake(%d, %d)", p.x, p.y)

	var dx, dy int
	switch p.direction {
	case api.MoveUp:
		dx = 0
		dy = 1
	case api.MoveDown:
		dx = 0
		dy = -1
	case api.MoveLeft:
		dx = 1
		dy = 0
	case api.MoveRight:
		dx = -1
		dy = 0
	}

	x := p.x
	y := p.y

	for i := p.size; i >= 0; i-- {
		board.SetCell(x, y, i)
		if x+dx < 0 || x+dx >= board.width {
			dx = 0
			dy = 1
		}
		if y+dy < 0 || y+dy >= board.height {
			dx = 1
			dy = 0
		}
		x += dx
		y += dy
	}
}

func (p *Player) Move(board *Board) error {
	var dx, dy int
	switch p.direction {
	case api.MoveLeft:
		dx = -1
		dy = 0
	case api.MoveRight:
		dx = 1
		dy = 0
	case api.MoveUp:
		dx = 0
		dy = -1
	case api.MoveDown:
		dx = 0
		dy = 1
	}

	nextX := p.x + dx
	nextY := p.y + dy

	if nextX < 0 || nextX == board.width || nextY < 0 || nextY == board.height {
		return fmt.Errorf("out of border")
	}
	if board.GetCell(nextX, nextY) > 0 {
		return fmt.Errorf("stamp snake")
	}
	if board.HitApple(nextX, nextY) {
		board.GenerateApple()
		p.size++
	}
	board.SetCell(nextX, nextY, p.size+1)
	p.x = nextX
	p.y = nextY
	return nil
}

func (p *Player) ChangeDirection(direction int) {
	// log.Printf("change direction: %d -> %d", p.direction, direction)
	// Do not turn around
	if p.direction == api.MoveDown && direction == api.MoveUp ||
		p.direction == api.MoveUp && direction == api.MoveDown ||
		p.direction == api.MoveLeft && direction == api.MoveRight ||
		p.direction == api.MoveRight && direction == api.MoveLeft {
		return
	}
	p.direction = direction
}
