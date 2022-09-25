package main

import (
	"math/rand"
)

type Player struct {
	ID    string
	X     int
	Y     int
	Theta int
	Object
}

func NewPlayer(id string, x, y int) *Player {
	return &Player{
		ID: id,
		X:  x,
		Y:  y,
		Object: Object{
			cood:  &Cood{Vector{X: x, Y: y}},
			theta: 0,
		},
	}
}

func (p *Player) Move() {
	v := 5
	cood := p.ConvertToWorld(&Cood{Vector{X: v, Y: 0}})
	p.cood.X = cood.X
	p.cood.Y = cood.Y

	rot := rand.Intn(10) - 5
	p.Rotate(rot)
}

func (p *Player) SetX(x int) {
	p.cood.X = x
}

func (p *Player) SetY(y int) {
	p.cood.Y = y
}

func (p *Player) SetDirection(theta int) {
	p.theta = theta
}
