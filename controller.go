package gameserver

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
		fmt.Printf("move: (%d, %d)\n", event.X, event.Y)
		p.Move(event.X, event.Y, event.Theta)
	default:
		fmt.Printf("unknown path: %s\n", msg.Path)
	}
}

func MatchMakingHandler(client *cambrian.WebSocketClient, engine interface{}) {
	ge := engine.(*system.GameEngine)
	switch ge.SceneMng.CurrentSceneID {
	case SceneMatchmaking:
		log.Printf("Scene: MatchMaking (%d)\n", len(ge.Clients))
		ge.AddClient(client)
		data, _ := json.Marshal(&api.SInitMessage{
			Status: api.GameStatusInit,
			ID:     client.ID(),
		})
		client.Send(data)
		if ge.ReachMaxClient() {
			ge.SceneMng.MoveScene(SceneIngame)

			players := make([]*Player, len(ge.Clients))
			for i, c := range ge.Clients {
				players[i] = NewPlayer(c, c.Stream(), 0, 0)
			}
			ingame = NewGame(1280, 960, players)
			go ingame.Run()
		} else {
			data := &api.EventResponse{
				Status: api.GameStatusWaiting,
			}

			bytes, _ := json.Marshal(&data)
			client.Send(bytes)
		}
	case SceneIngame:
		log.Printf("Scene: Ingame\n")
		ge.DeleteClient(client.ID())

		data := &api.EventResponse{
			Status: api.GameStatusError,
		}

		bytes, _ := json.Marshal(&data)
		client.Send(bytes)
	}
}

func DisconnectHandler(client *cambrian.WebSocketClient, engine interface{}) {
	fmt.Println("Disonnect!!")
	ge := engine.(*system.GameEngine)
	ge.DeleteClient(client.ID())
}

func PublishStatus(req cambrian.Request) {
	if ingame != nil && ingame.IsStart() {
		fmt.Println("--- tick!!!")
		players := make([]*Player, 0, len(ingame.players))
		for _, p := range ingame.players {
			players = append(players, p)
		}

		for _, player := range ingame.players {
			err := player.Send(api.GameStatusOK, players)

			if err != nil {
				player.Status = PlayerDead
				player.Finish()
				delete(ingame.players, player.ID())
			}
		}
	}
}
