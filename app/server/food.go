package server

import (
	"time"

	"jmhart.dev/htmx-argio/physics"
)

type Food struct {
	ID       string
	Pos      physics.Vec2f
	Color    string
	Consumed bool
	Collider physics.Collider
}

func (f *Food) New() {
	f.Consumed = false
	f.Collider.X = f.Pos.X
	f.Collider.Y = f.Pos.Y

	f.Collider.Width = 50
	f.Collider.Height = 50
}

func (f *Food) Consume(player *Player, server *Server) {
	player.Size += 5
	f.Consumed = true

	go func(f *Food) {
		time.Sleep(time.Minute * 1)
		f.Consumed = false
	}(f)
}
