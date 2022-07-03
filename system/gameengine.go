package system

import (
	"math/rand"
	"time"
)

type GameEngine struct {
	maxClient int
	Clients   []Client
	SceneMng  *SceneManager
}

func NewGameEngine(max int) *GameEngine {
	rand.Seed(time.Now().Unix())
	clients := make([]Client, 0)
	mng := NewSceneManager()
	return &GameEngine{
		maxClient: max,
		Clients:   clients,
		SceneMng:  mng,
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

func (ge *GameEngine) ReachMaxClient() bool {
	return len(ge.Clients) >= ge.maxClient
}
