package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xade-game/game-server/api"
	"github.com/xade-game/game-server/system"
)

const (
	PlayerNum = 2
)

func ingameHandler(mng *system.SceneManager, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	client := system.NewWebClient(c)
	log.Printf("Connect new websocket")
	go client.Run()
	client.AddObserver(mng)
	client.Notify(system.EventClientConnect)
}

func main() {
	var (
		addr string
	)

	flag.StringVar(&addr, "addr", ":8080", "http service address")
	flag.Parse()

	var err error

	ge := system.NewGameEngine(PlayerNum)
	ge.SceneMng.AddHandler(system.EventClientConnect, SceneMatchmaking, func(args interface{}) {
		log.Printf("Scene: MatchMaking (%d)\n", len(ge.Clients))
		ta := args.(system.TriggerArgument)
		ge.AddClient(ta.Client)
		ta.Client.Send([]byte(fmt.Sprintf("{\"status\":%d, \"id\": \"%s\"}", api.GameStatusInit, ta.Client.ID())))
		if ge.ReachMaxClient() {
			ge.SceneMng.MoveScene(SceneIngame)

			if err != nil {
				log.Fatalf("Agones SDK: Failed to Allocate: %v", err)
			}

			players := make([]*Player, len(ge.Clients))
			for i, c := range ge.Clients {
				players[i] = NewPlayer(c, c.Stream(), 0, 0)
			}
			ingame := NewGame(1280, 960, players)
			go ingame.Run()
		} else {
			data := &api.EventResponse{
				Status: api.GameStatusWaiting,
			}

			bytes, _ := json.Marshal(&data)
			ta.Client.Send(bytes)
		}
	})

	ge.SceneMng.AddHandler(system.EventClientConnect, SceneIngame, func(args interface{}) {
		log.Printf("Scene: Ingame\n")
		ta := args.(system.TriggerArgument)
		ge.DeleteClient(ta.Client.ID())

		data := &api.EventResponse{
			Status: api.GameStatusError,
		}

		bytes, _ := json.Marshal(&data)
		ta.Client.Send(bytes)
	})

	ge.SceneMng.AddHandler(system.EventClientFinish, SceneIngame, func(args interface{}) {
		log.Printf("Trigger: EventClientFinish\n")
		ta := args.(system.TriggerArgument)
		ge.DeleteClient(ta.Client.ID())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ingameHandler(ge.SceneMng, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
