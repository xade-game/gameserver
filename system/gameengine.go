package system

import (
	"math/rand"
	"time"
)

type Game interface {
	Update()
}

type GameEngine struct {
	Clients  []Client
	SceneMng *SceneManager
	game     Game
}

func NewGameEngine() *GameEngine {
	rand.Seed(time.Now().Unix())
	clients := make([]Client, 0)
	mng := NewSceneManager()
	return &GameEngine{
		Clients:  clients,
		SceneMng: mng,
	}
}

func (ge *GameEngine) AddClient(c Client) {
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

func (ge *GameEngine) ClientNum() int {
	return len(ge.Clients)
}

func (ge *GameEngine) SetGame(game Game) {
	ge.game = game
}

func (ge *GameEngine) Update() {
	if ge.game != nil {
		ge.game.Update()
	}
}
