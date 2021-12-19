package main

import (
	"fmt"
	"time"
)

func Exec(sm *GameStateMachine, f chan int) {
	for {
		time.Sleep(time.Duration(100) * time.Millisecond)
		if sm.State != initialized {
			break
		}
		c := RandomClient(sm.input)
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
		fmt.Println("--------------- Game Start!! --------------------")
		sm := NewGameStateMachine(f)
		go sm.ClientEventHandler()
		go sm.Run()
		Exec(sm, f)
		fmt.Println("--------------- Game Finish!! --------------------")
	}
}
