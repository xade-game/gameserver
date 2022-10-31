package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/xade-game/gameserver/api"
	"github.com/xade-game/gameserver/system"
)

var (
	ErrClientNotFound = errors.New("client not found")
)

const (
	PlayerNum  = 2
	GameWidth  = 1024
	GameHeight = 1024
)

type Game struct {
	width   int
	height  int
	engine  *system.GameEngine
	players map[string]*Player
}

func NewGame(w, h int, engine *system.GameEngine) *Game {
	players := make(map[string]*Player)

	for _, client := range engine.Clients {
		x := rand.Intn(GameWidth)
		y := rand.Intn(GameHeight)
		theta := rand.Intn(360)
		players[client.ID()] = NewPlayer(client, x, y, theta)
	}
	return &Game{
		width:   w,
		height:  h,
		engine:  engine,
		players: players,
	}
}

func (g *Game) SendStart() {
	presp := make([]api.PlayerResponse, len(g.players))
	resp := api.ResponseBody{
		Board:   make([]int, 0),
		Width:   GameWidth,
		Height:  GameHeight,
		Players: presp,
	}
	i := 0
	for _, p := range g.players {
		resp.Players[i] = p.ToResponse()
		i++
	}
	for _, p := range g.players {
		data := &api.EventResponse{
			Status: api.GameStatusOK,
			Body:   resp,
		}
		bytes, _ := json.Marshal(&data)
		p.Send(bytes)
	}
}

func (g *Game) Update() {
	fmt.Println("update")
	g.SendStart()
}

func (g *Game) UpdateClientStatus(req api.EventRequest) error {
	c, found := g.players[req.UUID]
	if !found {
		return ErrClientNotFound
	}

	c.X = req.X
	c.Y = req.Y
	c.Theta = req.Theta

	return nil
}

type Player struct {
	client system.Client
	X      int
	Y      int
	Theta  int
}

func NewPlayer(c system.Client, x, y, theta int) *Player {
	return &Player{
		client: c,
		X:      x,
		Y:      y,
		Theta:  theta,
	}
}

func (p *Player) ToResponse() api.PlayerResponse {
	return api.PlayerResponse{
		ID:        p.client.ID(),
		X:         p.X,
		Y:         p.Y,
		Direction: p.Theta,
	}
}

func (p *Player) Send(data []byte) error {
	return p.client.Send(data)
}
