package main

import (
	"github.com/xade-game/gameserver/api"
	"github.com/xade-game/gameserver/system"
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
	engine  *system.GameEngine
}

func NewGame(w, h int, engine *system.GameEngine) *Game {
	playerMap := make(map[string]*Player)

	board := NewBoard(w, h)
	board.GenerateApple()

	for _, c := range ge.Clients {
		player := NewPlayer(c, c.Stream(), w, h)
		playerMap[player.ID()] = player
		board.SetCell(player.x, player.y, 1)

		player.GenerateSnake(board)
	}

	return &Game{
		width:   w,
		height:  h,
		players: playerMap,
		status:  -1,
		board:   board,
		engine:  engine,
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

func (g *Game) Update() {
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

func (g *Game) SendAll() {
	players := g.PlayerArray()
	for _, player := range players {
		player.Send(api.GameStatusOK, ingame.board, players)
	}
}
