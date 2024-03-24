package server

import (
	"bytes"
	"log"

	"github.com/gorilla/websocket"
	"jmhart.dev/htmx-argio/physics"
	"jmhart.dev/htmx-argio/utils"
)

const move_speed float64 = 0.1

type Controls struct {
	Up    bool
	Down  bool
	Left  bool
	Right bool
}

type Viewport struct {
	Pos    physics.Vec2f
	Width  int
	Height int
}

func (v *Viewport) Center() physics.Vec2f {
	center := physics.Vec2f{X: float64(v.Width), Y: float64(v.Height)}
	center.DivF(2)
	return center
}
func (v *Viewport) WorldOffset() physics.Vec2f {
	offset := v.Pos
	offset.MultF(-1)
	return offset
}
func (v *Viewport) PlayerAbsolute() physics.Vec2f {
	center := v.Center()
	player_pos := v.Pos
	player_pos.Add(center)
	return player_pos
}

type Player struct {
	ID       string
	Username string
	Color    string
	Size     int
	Ctl      Controls
	Viewport Viewport
	Conn     *websocket.Conn
	Vel      physics.Vec2f
	Pos      physics.Vec2f
	Collider physics.Collider
}

func (p *Player) New(conn *websocket.Conn, id string) {
	p.Conn = conn
	p.ID = id
	p.Size = 100
	p.Color = utils.GenerateRandomHexColor()
	p.Username = "really long username 1234"
	p.Ctl = Controls{}
}

func (p *Player) SendTemplate(s *Server, template string, data any) error {
	var buf bytes.Buffer
	writer := &buf
	err := s.Templates.ExecuteTemplate(writer, template, data)
	if err != nil {
		return err
	}
	p.Conn.WriteMessage(
		websocket.TextMessage,
		buf.Bytes())
	return nil
}

func (p *Player) update(delta int64) {

	if p.Ctl.Up {
		p.Vel.Y += move_speed * float64(delta)
	}
	if p.Ctl.Down {
		p.Vel.Y -= move_speed * float64(delta)
	}
	if p.Ctl.Left {
		p.Vel.X += move_speed * float64(delta)
	}
	if p.Ctl.Right {
		p.Vel.X -= move_speed * float64(delta)
	}

	p.Vel.MultF(0.80)

	p.Vel.LimitF(30)

	p.Viewport.Pos.Add(p.Vel)

	p.Pos = p.Viewport.PlayerAbsolute()
	p.Pos.SubF(float64(p.Size / 2))

	p.Collider.X = p.Pos.X
	p.Collider.Y = p.Pos.Y

	p.Collider.Width = float64(p.Size)
	p.Collider.Height = float64(p.Size)
}

func (p *Player) sendPlayer(s *Server) {
	p.SendTemplate(s, "self.tmpl.html", p)
}

func (p *Player) sendPostion(s *Server) {
	err := p.SendTemplate(s, "self-position.tmpl.html", p)
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}
