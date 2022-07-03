package main

import (
	"time"

	"github.com/xade-game/game-server/api"
)

const (
	SceneMatchmaking = iota
	SceneIngame
)

type Game struct {
	width   int
	height  int
	players map[string]*Player
}

func NewGame(w, h int, players []*Player) *Game {
	playerMap := make(map[string]*Player)

	for _, p := range players {
		playerMap[p.ID()] = p
	}

	return &Game{
		width:   w,
		height:  h,
		players: playerMap,
	}
}

func (g *Game) Run() {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	players := make([]*Player, 0, len(g.players))
	for _, p := range g.players {
		players = append(players, p)
	}

	for range t.C {
		for _, player := range g.players {
			err := player.Send(api.GameStatusOK, players)

			if err != nil {
				player.Status = PlayerDead
				player.Finish()
				delete(g.players, player.ID())
			}
		}
	}
}
