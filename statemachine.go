package main

import (
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
		client.recv <- data
	}
}

type ClientData struct {
	Id     int    `json:"id"`
	Status string `json:"status"`
}

func (sm *GameStateMachine) ClientEventHandler() {
	for data := range sm.input {
		client := &ClientData{}
		json.Unmarshal(data, client)
		sm.UpdateClientState(client.Id, client.Status)
	}
}

func (sm *GameStateMachine) UpdateClientState(id int, status string) {
	for _, c := range sm.clients {
		if c.id == id {
			c.status = dead
			return
		}
	}
}

func (sm *GameStateMachine) Run() {
	go Tick(100, sm.tick)

	fmt.Println("Phase: Init")
	for range sm.tick {
		switch sm.State {
		case initialized:
			if sm.isOpened() {
				sm.State = opened
				sm.timer = 0
				fmt.Println("\nPhase: Opened")
				sm.Broadcast([]byte("opended"))
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
	return len(sm.clients) >= 10 || sm.timer > 20
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
