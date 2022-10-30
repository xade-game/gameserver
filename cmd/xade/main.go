package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xade-game/gameserver/cambrian"
	"github.com/xade-game/gameserver/system"
)

var ge *system.GameEngine

// var ingame *Game

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

	cmbr := cambrian.New()
	cmbr.RegisterWebsocketConnect(connectHandler)
	cmbr.RegisterWebsocketMessage(msgHandler)
	cmbr.RegisterWebsocketDisconnect(disconnectHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ingameHandler(ge.SceneMng, cmbr, w, r)
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
