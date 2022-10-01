package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/xade-game/gameserver/api"
	"github.com/xade-game/gameserver/cambrian"
	"github.com/xade-game/gameserver/system"
)

var ingame *Game

func RouteHandler(req cambrian.Request, engine interface{}) {
	var msg api.Message

	err := json.Unmarshal(req.Body(), &msg)
	if err != nil {
		log.Print(err)
	}

	switch msg.Path {
	case "move":
		var event api.EventRequest
		p, ok := ingame.FindPlayerById(req.ID())
		if !ok {
			return
		}

		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			log.Print(err)
		}
		p.ChangeDirection(event.Key)
	default:
		fmt.Printf("unknown path: '%s'\n", msg.Path)
	}
}

func MatchMakingHandler(client *cambrian.WebSocketClient, engine interface{}) {
	ge := engine.(*system.GameEngine)
	ge.AddClient(client)
	data, _ := json.Marshal(&api.SInitMessage{
		Status: api.GameStatusInit,
		ID:     client.ID(),
	})
	client.Send(data)
	if ge.ClientNum() >= PlayerNum {
		players := make([]*Player, 0, len(ge.Clients))
		for _, c := range ge.Clients {
			players = append(players, NewPlayer(c, c.Stream(), GameCellWidth, GameCellHeight))
		}
		ingame = NewGame(GameCellWidth, GameCellHeight, players)
		ingame.Start()

		// init game

		for _, player := range players {
			player.GenerateSnake(ingame.board)
		}

		for _, player := range players {
			player.Send(api.GameStatusOK, ingame.board, players)
		}
	} else {
		data := &api.EventResponse{
			Status: api.GameStatusWaiting,
		}

		bytes, _ := json.Marshal(&data)
		client.Send(bytes)
	}
}

func DisconnectHandler(client *cambrian.WebSocketClient, engine interface{}) {
	fmt.Println("Disconnect!!")
	ge := engine.(*system.GameEngine)
	ge.DeleteClient(client.ID())
	if ge.ClientNum() == 0 {
		ingame.Stop()
	}
}

func PublishStatus(req cambrian.Request) {
	if ingame != nil && ingame.IsStart() {
		ingame.Run()
	}
}
