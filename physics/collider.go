package physics

type Collider struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// AABB
func (c *Collider) IsColliding(other *Collider) bool {
	return c.X < other.X+other.Width &&
		c.X+c.Width > other.X &&
		c.Y < other.Y+other.Height &&
		c.Y+c.Height > other.Y
}
