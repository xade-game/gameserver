package gamelogic

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"

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
				x := rand.Intn(GameCellWidth)
				y := rand.Intn(GameCellHeight)
				players[i] = NewPlayer(c, c.Stream(), x, y)
			}
			ingame = NewGame(GameCellWidth, GameCellHeight, players)

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
	case SceneIngame:
		log.Printf("Scene: Ingame\n")
		player, _ := ingame.GetPlayer(client.ID())
		player.Send(api.GameStatusError, ingame.board, ingame.PlayerArray())
	}
}

func DisconnectHandler(client *cambrian.WebSocketClient, engine interface{}) {
	fmt.Println("Disonnect!!")
	ge := engine.(*system.GameEngine)
	ge.DeleteClient(client.ID())
}

func PublishStatus(req cambrian.Request) {
	if ingame != nil && ingame.IsStart() {
		log.Println("--- tick!!!")
		players := ingame.PlayerArray()

		for _, player := range ingame.players {
			err := player.Send(api.GameStatusOK, ingame.board, players)

			if err != nil {
				player.Status = PlayerDead
				player.Finish()
				delete(ingame.players, player.ID())
			}
		}
	}
}
