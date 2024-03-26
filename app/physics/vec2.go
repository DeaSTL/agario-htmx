package physics

import "math"

type Vec2f struct {
	X float64
	Y float64
}

func (v *Vec2f) Add(other Vec2f) {
	v.X += other.X
	v.Y += other.Y
}

func (v *Vec2f) Sub(other Vec2f) {
	v.X -= other.X
	v.Y -= other.Y
}

func (v *Vec2f) Mult(other Vec2f) {
	v.X *= other.X
	v.Y *= other.Y
}

func (v *Vec2f) Div(other Vec2f) {
	v.X /= other.X
	v.Y /= other.Y
}

func (v *Vec2f) AddF(val float64) {
	v.Add(Vec2f{X: val, Y: val})
}
func (v *Vec2f) SubF(val float64) {
	v.Sub(Vec2f{X: val, Y: val})
}
func (v *Vec2f) MultF(val float64) {
	v.Mult(Vec2f{X: val, Y: val})
}
func (v *Vec2f) DivF(val float64) {
	v.Div(Vec2f{X: val, Y: val})
}

func (v *Vec2f) LimitF(val float64) {
	if v.X > 0 {
		v.X = math.Min(v.X, val)
	} else {
		v.X = math.Max(v.X, -val)
	}
	if v.Y > 0 {
		v.Y = math.Min(v.Y, val)
	} else {
		v.Y = math.Max(v.Y, -val)
	}
}

func (v *Vec2f) SetRad(angle float64) {
	x, y := v.X, v.Y

	cos_theta := math.Cos(angle)
	sin_theta := math.Sin(angle)

	rotated_x := x*cos_theta - y*sin_theta
	rotated_y := x*sin_theta + y*cos_theta

	v.X = rotated_x
	v.Y = rotated_y
}

func (v *Vec2f) SetDeg(angle float64) {
	v.SetRad((180 / angle) * math.Pi)
}
