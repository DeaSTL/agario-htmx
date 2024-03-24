package server

import (
	"bytes"

	"github.com/gorilla/websocket"
	"jmhart.dev/htmx-argio/physics"
	"jmhart.dev/htmx-argio/utils"
)

var move_speed float64 = 15

type Controls struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

type Player struct {
	ID    string
	Color string
	Size  int
	Ctl   Controls
	Conn  *websocket.Conn
	Vel   physics.Vec2f
	Pos   physics.Vec2f
}

func (p *Player) New(conn *websocket.Conn, id string) {
	p.Conn = conn
	p.ID = id
	p.Size = 100
	p.Color = utils.GenerateRandomHexColor()
	p.Ctl = Controls{}
}

func (p *Player) update() {

	if p.Ctl.Up {
		p.Vel.Y -= move_speed
	}
	if p.Ctl.Down {
		p.Vel.Y += move_speed
	}
	if p.Ctl.Left {
		p.Vel.X -= move_speed
	}
	if p.Ctl.Right {
		p.Vel.X += move_speed
	}

	p.Pos.Add(p.Vel)

	p.Vel.MultF(0.80)

	p.Vel.LimitF(30)
}

func (p *Player) sendPlayer(s *Server) {
	var buf bytes.Buffer
	writer := &buf
	s.Templates.ExecuteTemplate(writer, "self.tmpl.html", p)
	p.Conn.WriteMessage(
		websocket.TextMessage,
		buf.Bytes())
}

func (p *Player) sendPlayerPostion(s *Server) {
	var buf bytes.Buffer
	writer := &buf
	s.Templates.ExecuteTemplate(writer, "self-position.tmpl.html", p)
	p.Conn.WriteMessage(
		websocket.TextMessage,
		buf.Bytes())
}
