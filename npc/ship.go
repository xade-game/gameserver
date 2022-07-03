package main

const (
	r = 30.0
)

type Ship struct {
	Position *Cood
	Width    int
	Height   int
	Object
}

func NewShip(x, y, theta int) *Ship {
	return &Ship{
		Object: Object{
			cood:  &Cood{Vector{X: x, Y: y}},
			theta: theta,
		},
		Width:  2*r + 30,
		Height: 2*r + 30,
	}
}
