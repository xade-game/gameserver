package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func Tick(msec int, t chan int) {
	for {
		time.Sleep(time.Duration(msec) * time.Millisecond)
		t <- 1
	}
}

type GameState int

const (
	initialized = iota
	opened
	closed
)

type GameStateMachine struct {
	State   GameState
	tick    chan int
	finish  chan int
	input   chan []byte
	clients []*Client
	timer   int
}

func NewGameStateMachine(f chan int) *GameStateMachine {
	t := make(chan int)
	i := make(chan []byte, 30)
	return &GameStateMachine{
		State:   initialized,
		tick:    t,
		finish:  f,
		input:   i,
		clients: []*Client{},
		timer:   0,
	}
}

func (sm *GameStateMachine) RegisterClient(client *Client) {
	if sm.State == initialized {
		client.status = registered
		sm.clients = append(sm.clients, client)
	}
}

func (sm *GameStateMachine) Broadcast(data []byte) {
	for _, client := range sm.clients {
		client.SendData(data)
	}
}

type CommandData struct {
	ClientId int    `json:"client_id"`
	Command  string `json:"command"`
	Status   string `json:"status"`
}

type GameStateData struct {
	Status string `json:"status"`
}

func (sm *GameStateMachine) ClientEventHandler(ctx context.Context) {
	for {
		select {
		case data := <-sm.input:
			cmd := &CommandData{}
			json.Unmarshal(data, cmd)
			sm.UpdateByCommand(cmd)
		case <-ctx.Done():
			return
		}
	}
}

func (sm *GameStateMachine) UpdateByCommand(cmd *CommandData) {
	switch cmd.Command {
	case "update":
		for _, c := range sm.clients {
			if c.id == cmd.ClientId {
				c.status = dead
				fmt.Printf(".")
				return
			}
		}
	}
}

func (sm *GameStateMachine) Run(ctx context.Context) {
	go Tick(100, sm.tick)

	fmt.Println("Phase: Init")
	for range sm.tick {
		switch sm.State {
		case initialized:
			if sm.isOpened() {
				sm.State = opened
				sm.timer = 0
				fmt.Println("\nPhase: Opened")
				data := &GameStateData{
					Status: "opened",
				}
				jsonData, _ := json.Marshal(data)
				sm.Broadcast(jsonData)
			}
		case opened:
			if sm.isClosed() {
				sm.State = closed
				sm.timer = 0
				fmt.Println("\nPhase: Closed")
			}
		case closed:
			for _, c := range sm.clients {
				if c.status == started {
					fmt.Printf("client(%d) is winner\n", c.id)
				}
			}
			sm.finish <- 1
			return
		}
		sm.timer += 1
	}
}

func (sm *GameStateMachine) isOpened() bool {
	return len(sm.clients) >= 10
}

func (sm *GameStateMachine) isClosed() bool {
	survivor := 0
	for _, c := range sm.clients {
		if c.status == started {
			survivor += 1
		}
	}

	return survivor <= 1
}
