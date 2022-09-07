package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xade-game/gameserver/api"
)

type EventRequest struct {
	UUID      string `json:"uuid"`
	Eventtype int    `json:"eventtype"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Theta     int    `json:"theta"`
}

type Message struct {
	UUID string `json:"uuid"`
	Path string `json:"path"`
	Body []byte `json:"body"`
}

var status = 0

func run(conn *websocket.Conn) {
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Close!! %v", err)
			conn.Close()
			return
		}
		fmt.Printf("recv: (%d) %s\n", mt, message)
		var msg api.EventResponse
		json.Unmarshal(message, &msg)
		status = msg.Status
	}
}

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		run(conn)
	}()
	defer conn.Close()

	conn2, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Millisecond)

	go func() {
		run(conn2)
	}()
	defer conn2.Close()

	i := 0

	for {
		if i >= 10 {
			return
		}
		if status == 1 {
			event, _ := json.Marshal(&EventRequest{
				UUID:      "hogehoge",
				Eventtype: 0,
				X:         100,
				Y:         20,
				Theta:     120,
			})
			data, _ := json.Marshal(&Message{
				Path: "move",
				UUID: "hoge",
				Body: event,
			})
			err = conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Fatal(err)
			}
			i++
		}
		time.Sleep(50 * time.Millisecond)
	}
}
