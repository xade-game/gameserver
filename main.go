package main

import (
	"context"
	"fmt"
	"time"
)

func Exec(ctx context.Context, sm *GameStateMachine, f chan int) {
	go sm.Tick()

	for {
		time.Sleep(time.Duration(100) * time.Millisecond)
		if sm.State != initialized {
			break
		}
		c := RandomClient(ctx, sm.input)
		if c != nil {
			sm.RegisterClient(c)
		}
	}

	for range f {
		return
	}
}

func main() {
	f := make(chan int)
	for {
		fmt.Println("--------------- Game Created!! --------------------")
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		sm := NewGameStateMachine(ctx, f)

		Exec(ctx, sm, f)
		cancel()
		fmt.Println("--------------- Game Finish!! --------------------")
	}
}
