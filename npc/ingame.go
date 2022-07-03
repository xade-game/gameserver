package main

func NewIngameScene() *IngameScene {
	return &IngameScene{}
}

type IngameScene struct {
}

func (s *IngameScene) Start() {
}

func (s *IngameScene) Update() (SceneType, error) {
	player := game.Player
	player.Move()
	game.conn.SendData(player.cood.X, player.cood.Y, player.theta)

	return SceneType("ingame"), nil
}

func (s *IngameScene) Finish() {
	game.conn.Close()
}
