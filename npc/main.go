package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/myoan/snake/api"
)

var game *Game

type Status int

type Game struct {
	id       string
	sceneMng *SceneManager
	conn     *Conn
	Status   int
	Players  []*Player
	Player   *Player
}

func (g *Game) Update() error {
	return g.sceneMng.Update()
}

func main() {
	var addr string

	flag.StringVar(&addr, "addr", ":8080", "address")
	flag.Parse()

	game = &Game{
		Status: StatusInit,
	}

	mng := NewSceneManager()
	game.conn = NewConn()
	game.sceneMng = mng
	game.sceneMng.AddScene("menu", NewMenuScene(addr))
	game.sceneMng.AddScene("matchmaking", NewMatchmakingScene())
	game.sceneMng.AddScene("ingame", NewIngameScene())
	game.sceneMng.SetInitialScene("menu")

	game.conn.AddHandler(api.GameStatusInit, func(message []byte) error {
		var resp api.InitResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			fmt.Printf("err %v\n", err)
			return err
		}
		fmt.Printf("resp: %+v\n", resp)
		game.id = resp.ID
		game.Player = NewPlayer(resp.ID, 0, 0)
		return nil
	})
	game.conn.AddHandler(api.GameStatusOK, func(message []byte) error {
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}
		if game.Status == StatusInit || game.Status == StatusWait {
			game.Status = StatusStart
			for _, p := range resp.Body.Players {
				if p.ID == game.id {
					game.Player.SetX(p.X)
					game.Player.SetY(p.Y)
					game.Player.SetDirection(p.Direction)
				} else {
					player := NewPlayer(p.ID, p.X, p.Y)
					game.Players = append(game.Players, player)
				}
			}
		} else if game.Status == StatusStart {
			for _, p := range resp.Body.Players {
				if p.ID != game.id {
					for _, player := range game.Players {
						player.SetX(p.X)
						player.SetY(p.Y)
						player.SetDirection(p.Direction)
					}
				}
			}
		}
		return nil
	})
	game.conn.AddHandler(api.GameStatusError, func(message []byte) error {
		game.Status = StatusDrop
		var resp api.EventResponse
		err := json.Unmarshal(message, &resp)
		if err != nil {
			return err
		}

		return fmt.Errorf("error")
	})
	game.conn.AddHandler(api.GameStatusWaiting, func(message []byte) error {
		game.Status = StatusWait
		return nil
	})

	t := time.NewTicker(100 * time.Millisecond)
	for range t.C {
		err := game.Update()
		if err != nil {
			panic(err)
		}
	}
}
