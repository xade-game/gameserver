package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xade-game/gameserver/cambrian"
	"github.com/xade-game/gameserver/system"
)

const (
	PlayerNum = 2
)

var ge *system.GameEngine

func ingameHandler(mng *system.SceneManager, cmbr *cambrian.Cambrian, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	cmbr.AddWebsocketClient(c, ge)
}

func main() {
	var addr string

	flag.StringVar(&addr, "addr", ":8080", "http service address")
	flag.Parse()

	ge = system.NewGameEngine()

	ge.SceneMng.AddHandler(system.EventClientFinish, SceneIngame, func(args interface{}) {
		ta := args.(system.TriggerArgument)
		ge.DeleteClient(ta.Client.ID())
	})

	cmbr := cambrian.New()
	cmbr.RegisterWebsocketConnect(MatchMakingHandler)
	cmbr.RegisterWebsocketMessage(RouteHandler)
	cmbr.RegisterWebsocketDisconnect(DisconnectHandler)
	cmbr.RegisterPeriodic(100*time.Millisecond, PublishStatus)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ingameHandler(ge.SceneMng, cmbr, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
