package gamelogic

import (
	"math/rand"
	"time"

	"github.com/xade-game/gameserver/api"
)

const (
	SceneMatchmaking = iota
	SceneIngame
)

type Game struct {
	width   int
	height  int
	players map[string]*Player
	status  int
	board   *Board
}

func NewGame(w, h int, players []*Player) *Game {
	playerMap := make(map[string]*Player)

	board := NewBoard(w, h)
	board.GenerateApple()

	for _, p := range players {
		playerMap[p.ID()] = p
		board.SetCell(p.x, p.y, 1)
	}

	return &Game{
		width:   w,
		height:  h,
		players: playerMap,
		status:  -1,
		board:   board,
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

func (g *Game) DrawBoard() {
	for _, p := range g.players {
		g.board.SetCell(p.x, p.y, 3)
	}
}

func (g *Game) RefreshUser() {
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
			err := player.Send(api.GameStatusOK, g.board, players)

			// player.Move()

			if err != nil {
				player.Status = PlayerDead
				player.Finish()
				delete(g.players, player.ID())
			}
		}
	}
}

type Board struct {
	board  [][]int
	width  int
	height int
}

func NewBoard(w, h int) *Board {
	board := make([][]int, h)
	for i := range board {
		board[i] = make([]int, w)
	}
	return &Board{
		board:  board,
		width:  w,
		height: h,
	}
}

func (b *Board) Reset() {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			if b.board[y][x] > 0 {
				b.board[y][x] = 0
			}
		}
	}
}

func (b *Board) GenerateApple() {
	for {
		x := rand.Intn(b.width)
		y := rand.Intn(b.height)

		if b.GetCell(x, y) == 0 {
			b.SetCell(x, y, -1)
			return
		}
	}
}

func (b *Board) Update() {
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			if b.board[i][j] > 0 {
				b.board[i][j]--
			}
		}
	}
}

func (b *Board) HitApple(x, y int) bool {
	return b.board[y][x] == -1
}

func (b *Board) GetCell(x, y int) int {
	return b.board[y][x]
}

func (b *Board) SetCell(x, y, data int) {
	b.board[y][x] = data
}

func (b *Board) ToArray() []int {
	ret := make([]int, b.width*b.height)

	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			ret[y*b.width+x] = b.board[y][x]
		}
	}
	return ret
}