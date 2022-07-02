package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/myoan/snake/api"
)

const (
	PlayerNum = 2
)

const (
	EventClientConnect = iota
	EventClientFinish
	EventClientRestart
)

type Observer interface {
	Update(data interface{}) error
}

type TriggerArgument struct {
	EventType int
	Client    Client
}

func ingameHandler(mng *SceneManager, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	stream := make(chan []byte)
	obs := make([]Observer, 0)
	client := &WebClient{
		uuid:      uuid.NewString(),
		stream:    stream,
		conn:      c,
		observers: obs,
	}

	log.Printf("Connect new websocket")
	go client.Run(stream)
	client.AddObserver(mng)
	client.Notify(EventClientConnect)
}

func main() {
	var (
		addr string
	)

	flag.StringVar(&addr, "addr", ":8080", "http service address")
	flag.Parse()

	var err error

	ge := NewGameEngine()
	ge.SceneMng.AddHandler(EventClientConnect, SceneMatchmaking, func(args interface{}) {
		log.Printf("Scene: MatchMaking (%d)\n", len(ge.Clients))
		ta := args.(TriggerArgument)
		ge.AddClient(ta.Client)
		ta.Client.Send([]byte(fmt.Sprintf("{\"status\":%d, \"id\": \"%s\"}", api.GameStatusInit, ta.Client.ID())))
		if ge.ReachMaxClient() {
			ge.SceneMng.MoveScene(SceneIngame)

			if err != nil {
				log.Fatalf("Agones SDK: Failed to Allocate: %v", err)
			}

			ge.ExecuteIngame()
		} else {
			data := &api.EventResponse{
				Status: api.GameStatusWaiting,
			}

			bytes, _ := json.Marshal(&data)
			ta.Client.Send(bytes)
		}
	})

	ge.SceneMng.AddHandler(EventClientConnect, SceneIngame, func(args interface{}) {
		log.Printf("Scene: Ingame\n")
		ta := args.(TriggerArgument)
		ge.DeleteClient(ta.Client.ID())

		data := &api.EventResponse{
			Status: api.GameStatusError,
		}

		bytes, _ := json.Marshal(&data)
		ta.Client.Send(bytes)
	})

	ge.SceneMng.AddHandler(EventClientFinish, SceneIngame, func(args interface{}) {
		log.Printf("Trigger: EventClientFinish\n")
		ta := args.(TriggerArgument)
		ge.DeleteClient(ta.Client.ID())
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ingameHandler(ge.SceneMng, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
