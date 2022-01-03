package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bykof/stateful"
)

type GameState struct {
	state   stateful.State
	tick    chan int
	done    chan int
	input   chan []byte
	clients []*Client
	timer   int
}

func (gs GameState) State() stateful.State {
	return gs.state
}

func (gs *GameState) SetState(state stateful.State) error {
	gs.state = state
	return nil
}

func (gs *GameState) Start(targs stateful.TransitionArguments) (stateful.State, error) {
	_, ok := targs.(GameArguments)
	if !ok {
		return nil, errors.New("could not parse GameArguments")
	}

	if gs.isOpened() {
		fmt.Println("state init -> opened")
		gs.timer = 0
		data := &GameStateData{
			Status: "opened",
		}
		jsonData, _ := json.Marshal(data)
		gs.Broadcast(jsonData)
		return Opened, nil
	}
	return Init, nil
}

func (gs *GameState) Finish(targs stateful.TransitionArguments) (stateful.State, error) {
	_, ok := targs.(GameArguments)
	if !ok {
		return nil, errors.New("could not parse GameArguments")
	}

	if gs.isClosed() {
		fmt.Println("\nstate opened -> closed")
		gs.timer = 0
		gs.done <- 1
		return Closed, nil
	}
	return Opened, nil
}

func (gs *GameState) RegisterClient(client *Client) {
	client.status = registered
	gs.clients = append(gs.clients, client)
}

func (gs *GameState) Broadcast(data []byte) {
	for _, client := range gs.clients {
		client.SendData(data)
	}
}

func (gs *GameState) isOpened() bool {
	return len(gs.clients) >= 10
}

func (gs *GameState) isClosed() bool {
	survivor := 0
	for _, c := range gs.clients {
		if c.status == started {
			survivor += 1
		}
	}

	return survivor <= 1
}
