package main

import (
	"context"
	"fmt"
)

func main() {
	f := make(chan int)
	for {
		fmt.Println("--------------- Game Created!! --------------------")
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		game := NewGameEngine(f)
		game.Run(ctx)
		cancel()

		fmt.Println("--------------- Game Finish!! --------------------")
	}
}
