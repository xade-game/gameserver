package main

import (
	"encoding/json"
	"fmt"

	"github.com/xade-game/gameserver/api"
	"github.com/xade-game/gameserver/cambrian"
	"github.com/xade-game/gameserver/system"
)

func connectHandler(client *cambrian.WebSocketClient, engine interface{}) {
	ge := engine.(*system.GameEngine)
	ge.AddClient(client)
	data, _ := json.Marshal(&api.SInitMessage{
		Status: api.GameStatusInit,
		ID:     client.ID(),
	})
	client.Send(data)
	if ge.ClientNum() >= PlayerNum {
		ingame = NewGame(GameWidth, GameHeight, ge)
		ingame.SendStart()
	} else {
		data := &api.EventResponse{
			Status: api.GameStatusWaiting,
		}

		bytes, _ := json.Marshal(&data)
		client.Send(bytes)
	}
}

func disconnectHandler(c *cambrian.WebSocketClient, obj interface{}) {
	fmt.Println("Disconnected")
}

func msgHandler(req cambrian.Request, obj interface{}) {
	body := req.Body()
	fmt.Printf("body: %s\n", body)

	var r api.EventRequest
	json.Unmarshal(body, &r)
	ingame.UpdateClientStatus(r)

}

func publishStatus(req cambrian.Request) {
	if ingame != nil {
		ingame.Update()
	}
}
