package main

import "math"

type Vector struct {
	X, Y int
}

func (v *Vector) Length() float64 {
	return math.Sqrt(float64(v.X*v.X + v.Y*v.Y))
}

func (v *Vector) Distance(dst *Vector) float64 {
	return math.Sqrt(float64(v.X*dst.X + v.Y*dst.Y))
}

func (v *Vector) Angle(dst *Vector) float64 {
	cos := float64(v.X*dst.X+v.Y*dst.Y) / (v.Length() * dst.Length())
	return cos
}

type Cood struct {
	Vector
}

func (c *Cood) PositionFromDistance(rot, len int) *Cood {
	rad := toRadian(rot)
	dx := float64(len) * math.Cos(rad)
	dy := float64(len) * math.Sin(rad)

	return &Cood{Vector{X: c.X + int(dx), Y: c.Y - int(dy)}}
}

func toRadian(phai int) float64 {
	return float64(phai) * math.Pi / 180
}

type Object struct {
	cood  *Cood
	theta int
}

func (o *Object) Rotate(r int) {
	o.theta = (o.theta + r) % 360
}

func (o *Object) ConvertToWorld(target *Cood) *Cood {
	rad := toRadian(-o.theta)
	x := float64(o.cood.X) + math.Cos(rad)*float64(target.X) - math.Sin(rad)*float64(target.Y)
	y := float64(o.cood.Y) + math.Sin(rad)*float64(target.X) + math.Cos(rad)*float64(target.Y)

	return &Cood{Vector{X: int(x), Y: int(y)}}
}

func (o *Object) ConvertToLocal(target *Cood) *Cood {
	// 平行移動
	tmpx := target.X - o.cood.X
	tmpy := target.Y - o.cood.Y

	// 回転
	rad := toRadian(o.theta)
	sin := math.Sin(rad)
	cos := math.Cos(rad)
	x := cos*float64(tmpx) - sin*float64(tmpy)
	y := sin*float64(tmpx) + cos*float64(tmpy)

	return &Cood{Vector{X: int(x), Y: int(y)}}
}
