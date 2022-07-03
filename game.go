package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/myoan/snake/api"
)

type GameEngine struct {
	Clients  []Client
	SceneMng *SceneManager
}

func NewGameEngine() *GameEngine {
	rand.Seed(time.Now().Unix())
	clients := make([]Client, 0)
	mng := NewSceneManager()
	return &GameEngine{
		Clients:  clients,
		SceneMng: mng,
	}
}

func (ge *GameEngine) AddClient(c Client) {
	ge.Clients = append(ge.Clients, c)
}

func (ge *GameEngine) DeleteClient(cid string) {
	for i, c := range ge.Clients {
		if c.ID() == cid {
			ge.Clients = append(ge.Clients[:i], ge.Clients[i+1:]...)
			return
		}
	}
}

func (ge *GameEngine) ReachMaxClient() bool {
	return len(ge.Clients) >= PlayerNum
}

func (ge *GameEngine) ExecuteIngame() {
	players := make([]*Player, len(ge.Clients))
	for i, c := range ge.Clients {
		players[i] = NewPlayer(c, c.Stream(), 0, 0)
	}
	ingame := NewGame(1280, 960, players)
	go ingame.Run()
}

const (
	SceneMatchmaking = iota
	SceneIngame
)

type Client interface {
	ID() string
	Send(data []byte) error
	Close()
	Stream() chan []byte
}

type Scene struct {
	ID       int
	eventMap map[int]SceneHandler
}

func (s *Scene) AddEventHandler(eventType int, h SceneHandler) error {
	fn := s.eventMap[eventType]
	if fn != nil {
		return fmt.Errorf("scene ID:'%d' already exists", eventType)
	}
	s.eventMap[eventType] = h
	return nil
}

type SceneHandler func(interface{})

type SceneManager struct {
	// SceneID is current scene ID
	SceneID        int
	sceneMap       map[int]func(args interface{})
	Scenes         []*Scene
	defaultHandler SceneHandler
}

func (h SceneHandler) handler(i interface{}) {
	h(i)
}

func NewSceneManager() *SceneManager {
	m := make(map[int]func(interface{}))
	scenes := make([]*Scene, 0)
	return &SceneManager{
		SceneID:        SceneMatchmaking,
		sceneMap:       m,
		Scenes:         scenes,
		defaultHandler: func(args interface{}) {},
	}
}

// FindBySceneID is return scene by sceneID
// if not found, return error
func (mng *SceneManager) FindBySceneID(sceneID int) (*Scene, error) {
	for _, scene := range mng.Scenes {
		if scene.ID == sceneID {
			return scene, nil
		}
	}
	return nil, fmt.Errorf("scene ID:'%d' not found", sceneID)
}

// AddHandler set function which called when server is selected scene and selected event is occurred
// If selected scene or event is not found, SceneManagaer call default handler and it do nothing.
// If you want to change default handler, you can use SceneManager.DefaultHandler(f)
func (mng *SceneManager) AddHandler(eventType int, sceneID int, h SceneHandler) error {
	scene, err := mng.FindBySceneID(sceneID)
	if err != nil {
		mng.addScene(sceneID)
		scene, _ = mng.FindBySceneID(sceneID)
	}
	scene.AddEventHandler(eventType, h)
	return nil
}

// DefaultHandler set default handler which called when selected scene and selected event is not found
func (mng *SceneManager) DefaultHandler(h SceneHandler) {
	mng.defaultHandler = h
}

func (mng *SceneManager) Update(data interface{}) error {
	args := data.(TriggerArgument)

	scene, err := mng.FindBySceneID(mng.SceneID)
	if err != nil {
		return err
	}

	fn := scene.eventMap[args.EventType]
	if fn == nil {
		return fmt.Errorf("scene ID:'%d' not found", args.EventType)
	}
	fn(args)
	return nil
}

func (mng *SceneManager) MoveScene(sid int) {
	mng.SceneID = sid
	fmt.Printf("SceneID Change: %d\n", sid)
}

func (mng *SceneManager) addScene(sceneID int) {
	m := make(map[int]SceneHandler)
	scene := &Scene{
		ID:       sceneID,
		eventMap: m,
	}
	mng.Scenes = append(mng.Scenes, scene)
}

type Game struct {
	width   int
	height  int
	players map[string]*Player
}

func NewGame(w, h int, players []*Player) *Game {
	playerMap := make(map[string]*Player)

	for _, p := range players {
		playerMap[p.ID()] = p
	}

	return &Game{
		width:   w,
		height:  h,
		players: playerMap,
	}
}

func (g *Game) Run() {
	t := time.NewTicker(100 * time.Millisecond)
	defer t.Stop()

	players := make([]*Player, 0, len(g.players))
	for _, p := range g.players {
		players = append(players, p)
	}

	for range t.C {
		for _, player := range g.players {
			err := player.Send(api.GameStatusOK, players)

			if err != nil {
				player.Status = PlayerDead
				player.Finish()
				delete(g.players, player.ID())
			}
		}
	}
}
