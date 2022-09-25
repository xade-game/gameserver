package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xade-game/gameserver/cambrian"
	"github.com/xade-game/gameserver/gamelogic"
	"github.com/xade-game/gameserver/system"
)

const (
	PlayerNum = 2
)

var cmbr *cambrian.Cambrian
var ge *system.GameEngine

func ingameHandler(mng *system.SceneManager, w http.ResponseWriter, r *http.Request) {
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

	ge = system.NewGameEngine(PlayerNum)

	ge.SceneMng.AddHandler(system.EventClientFinish, gamelogic.SceneIngame, func(args interface{}) {
		log.Printf("Trigger: EventClientFinish\n")
		ta := args.(system.TriggerArgument)
		ge.DeleteClient(ta.Client.ID())
	})

	cmbr = cambrian.New()
	cmbr.RegisterWebsocketConnect(gamelogic.MatchMakingHandler)
	cmbr.RegisterWebsocketMessage(gamelogic.RouteHandler)
	cmbr.RegisterWebsocketDisconnect(gamelogic.DisconnectHandler)
	cmbr.RegisterPeriodic(100*time.Millisecond, gamelogic.PublishStatus)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ingameHandler(ge.SceneMng, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
