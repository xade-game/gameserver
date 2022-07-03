package main

import (
	"math/rand"
	"time"

	"github.com/myoan/snake/api"
	"github.com/xade-game/game-server/system"
)

type GameEngine struct {
	Clients  []system.Client
	SceneMng *system.SceneManager
}

func NewGameEngine() *GameEngine {
	rand.Seed(time.Now().Unix())
	clients := make([]system.Client, 0)
	mng := system.NewSceneManager(SceneMatchmaking)
	return &GameEngine{
		Clients:  clients,
		SceneMng: mng,
	}
}

func (ge *GameEngine) AddClient(c system.Client) {
	ge.Clients = append(ge.Clients, c)
}

func (ge *GameEngine) DeleteClient(cid string) {
	for i, c := range ge.Clients {
		if c.ID() == cid {
			ge.Clients = append(ge.Clients[:i], ge.Clients[i+1:]...)
			return
		}
	}
}

func (ge *GameEngine) ReachMaxClient() bool {
	return len(ge.Clients) >= PlayerNum
}

func (ge *GameEngine) ExecuteIngame() {
	players := make([]*Player, len(ge.Clients))
	for i, c := range ge.Clients {
		players[i] = NewPlayer(c, c.Stream(), 0, 0)
	}
	ingame := NewGame(1280, 960, players)
	go ingame.Run()
}

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
