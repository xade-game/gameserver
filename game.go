package gameserver

const (
	SceneMatchmaking = iota
	SceneIngame
)

type Game struct {
	width   int
	height  int
	players map[string]*Player
	status  int
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
		status:  -1,
	}
}

func (g *Game) Start() {
	g.status = 0
}

func (g *Game) Stop() {
	g.status = -1
}

func (g *Game) IsStart() bool {
	return g.status == 0
}

func (g *Game) FindPlayerById(id string) (*Player, bool) {
	p, found := g.players[id]
	return p, found
}

func (g *Game) RefreshUser() {
}
