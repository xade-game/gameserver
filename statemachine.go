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
	em      *EventManager
}

func NewGameStateMachine(ctx context.Context, f chan int) *GameStateMachine {
	t := make(chan int)
	i := make(chan []byte, 30)
	sm := &GameStateMachine{
		State:   initialized,
		tick:    t,
		finish:  f,
		input:   i,
		clients: []*Client{},
		timer:   0,
		em:      NewEventManager(),
	}

	sm.em.AddEventListener("init", sm.GameInitHandler)
	sm.em.AddEventListener("opened", sm.GameOpenedHandler)
	sm.em.AddEventListener("closed", sm.GameClosedHandler)
	sm.em.AddEventListener("update", sm.ClientEventHandler)

	go sm.em.Run(ctx)
	go sm.SetClientEventDispatcher(ctx)

	return sm
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

func (sm *GameStateMachine) SetClientEventDispatcher(ctx context.Context) {
	for {
		select {
		case data := <-sm.input:
			sm.em.DispatchEvent("update", string(data))
		case <-ctx.Done():
			return
		}
	}
}

func (sm *GameStateMachine) ClientEventHandler(e *Event) {
	cmd := &CommandData{}
	json.Unmarshal([]byte(e.data), cmd)
	for _, c := range sm.clients {
		if c.id == cmd.ClientId {
			c.status = dead
			fmt.Printf(".")
			return
		}
	}
}

func (sm *GameStateMachine) GameInitHandler(e *Event) {
	if sm.isOpened() {
		fmt.Println("\nPhase: Opened")
		sm.State = opened
		sm.timer = 0
		data := &GameStateData{
			Status: "opened",
		}
		jsonData, _ := json.Marshal(data)
		sm.Broadcast(jsonData)
	}
}

func (sm *GameStateMachine) GameOpenedHandler(e *Event) {
	if sm.isClosed() {
		fmt.Println("\nPhase: Closed")
		sm.State = closed
		sm.timer = 0
	}
}

func (sm *GameStateMachine) GameClosedHandler(e *Event) {
	for _, c := range sm.clients {
		if c.status == started {
			fmt.Printf("client(%d) is winner\n", c.id)
		}
	}
	sm.finish <- 1
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

func (sm *GameStateMachine) Tick() {
	go Tick(100, sm.tick)
	for range sm.tick {
		sm.timer += 1
		switch sm.State {
		case initialized:
			sm.em.DispatchEvent("init")
		case opened:
			sm.em.DispatchEvent("opened")
		case closed:
			sm.em.DispatchEvent("closed")
		}
	}
}
