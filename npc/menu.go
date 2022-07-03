package main

import (
	"fmt"
)

func init() {
	fmt.Println("menu")
}

func NewMenuScene(addr string) *MenuScene {
	return &MenuScene{
		addr: addr,
	}
}

type MenuScene struct {
	addr string
}

func (s *MenuScene) Start() {}

func (s *MenuScene) Update() (SceneType, error) {
	go game.conn.Connect(s.addr)
	return SceneType("matchmaking"), nil
}
func (s *MenuScene) Finish() {}
