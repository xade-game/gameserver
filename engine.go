package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bykof/stateful"
)

/*
func Tick(msec int, t chan int) {
	for {
		time.Sleep(time.Duration(msec) * time.Millisecond)
		t <- 1
	}
}
*/

type CommandData struct {
	ClientId int    `json:"client_id"`
	Command  string `json:"command"`
	Status   string `json:"status"`
}

type GameStateData struct {
	Status string `json:"status"`
}

const (
	Init   = stateful.DefaultState("init")
	Opened = stateful.DefaultState("opened")
	Closed = stateful.DefaultState("closed")
)

type GameEngine struct {
	gameState *GameState
	sm        *stateful.StateMachine
}

type GameArguments struct{}

func NewGameEngine(f chan int) *GameEngine {
	t := make(chan int)
	i := make(chan []byte, 30)
	gs := &GameState{
		state:   Init,
		tick:    t,
		finish:  f,
		input:   i,
		clients: []*Client{},
		timer:   0,
	}

	stateMachine := &stateful.StateMachine{
		StatefulObject: gs,
	}

	stateMachine.AddTransition(
		gs.Start,
		stateful.States{Init},
		stateful.States{Opened},
	)
	stateMachine.AddTransition(
		gs.Finish,
		stateful.States{Opened},
		stateful.States{Closed},
	)

	return &GameEngine{
		gameState: gs,
		sm:        stateMachine,
	}
}

func (game *GameEngine) Run(ctx context.Context) {
	go game.SetClientEventDispatcher(ctx)

	for {
		time.Sleep(time.Duration(100) * time.Millisecond)
		if game.gameState.state != Init {
			break
		}
		c := RandomClient(ctx, game.gameState.input)
		if c != nil {
			game.gameState.RegisterClient(c)
			game.sm.Run(game.gameState.Start, stateful.TransitionArguments(GameArguments{}))
		}
	}

	for range game.gameState.finish {
		return
	}
}

func (game *GameEngine) SetClientEventDispatcher(ctx context.Context) {
	for {
		select {
		case data := <-game.gameState.input:
			cmd := &CommandData{}
			json.Unmarshal([]byte(data), cmd)
			for _, c := range game.gameState.clients {
				if c.id == cmd.ClientId {
					c.status = dead
					fmt.Printf(".")
					game.sm.Run(game.gameState.Finish, stateful.TransitionArguments(GameArguments{}))
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
