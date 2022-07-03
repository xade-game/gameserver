package main

import "time"

const (
	dotInterval = 30
)

func NewMatchmakingScene() *MatchmakingScene {
	return &MatchmakingScene{}
}

type MatchmakingScene struct{}

func (s *MatchmakingScene) Start() {
}

func (s *MatchmakingScene) Update() (SceneType, error) {
	for {
		time.Sleep(time.Second)
		if game.Status == StatusStart {
			return SceneType("ingame"), nil
		}
	}
}

func (s *MatchmakingScene) Finish() {}
