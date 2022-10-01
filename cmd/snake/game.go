package main

import (
	"github.com/xade-game/gameserver/api"
)

const (
	SceneMatchmaking = iota
	SceneIngame

	GameCellHeight = 40
	GameCellWidth  = 40
)

type Game struct {
	width   int
	height  int
	players map[string]*Player
	status  int
	board   *Board
}

func NewGame(w, h int, players []*Player) *Game {
	playerMap := make(map[string]*Player)

	board := NewBoard(w, h)
	board.GenerateApple()

	for _, p := range players {
		playerMap[p.ID()] = p
		board.SetCell(p.x, p.y, 1)
	}

	return &Game{
		width:   w,
		height:  h,
		players: playerMap,
		status:  -1,
		board:   board,
	}
}

func (g *Game) Start() {
	g.status = 0
}

func (g *Game) Stop() {
	g.status = -1
}

func (g *Game) IsStart() bool {
	return g.status == 0
}

func (g *Game) FindPlayerById(id string) (*Player, bool) {
	p, found := g.players[id]
	return p, found
}

func (g *Game) DrawBoard() {
	for _, p := range g.players {
		g.board.SetCell(p.x, p.y, 3)
	}
}

func (g *Game) RefreshUser() {
}

func (g *Game) Run() {
	players := g.PlayerArray()

	for _, player := range g.players {
		if err := player.Move(g.board); err != nil {
			player.Status = PlayerDead
			player.Finish()
			delete(g.players, player.ID())
		}
	}
	g.board.Update()
	for _, player := range g.players {
		err := player.Send(api.GameStatusOK, g.board, players)

		if err != nil {
			player.Status = PlayerDead
			player.Finish()
			delete(g.players, player.ID())
		}
	}
}

func (g *Game) GetPlayer(id string) (*Player, bool) {
	p, found := g.players[id]
	return p, found
}

func (g *Game) PlayerArray() []*Player {
	players := make([]*Player, 0, len(g.players))
	for _, p := range g.players {
		players = append(players, p)
	}
	return players
}
