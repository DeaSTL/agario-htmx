package server

import (
	"bytes"
	"log"
	"math"
	"sort"
	"sync"

	"github.com/gorilla/websocket"
	"jmhart.dev/htmx-argio/physics"
	"jmhart.dev/htmx-argio/utils"
)

var move_speed float64 = 0.1

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

type Leaderboard struct {
	Top     []*Player
	Current *Player
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
	ID           string
	Username     string
	Dead         bool
	Iinitialized bool
	Color        string
	Size         float64
	FontSize     float64
	EatPower     float64
	Ctl          Controls
	Viewport     Viewport
	Conn         *websocket.Conn
	Vel          physics.Vec2f
	Pos          physics.Vec2f
	Collider     physics.Collider
	sync.Mutex
}

type BySize []*Player

func (a BySize) Len() int           { return len(a) }
func (a BySize) Less(i, j int) bool { return a[i].Size > a[j].Size }
func (a BySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (p *Player) New(conn *websocket.Conn, id string) {
	p.Conn = conn
	p.ID = id
	p.Size = 100
	p.Color = utils.GenerateRandomHexColor()
	p.Username = "user_" + utils.GenID(6)
	p.Ctl = Controls{
		Left:  false,
		Up:    false,
		Right: false,
		Down:  false,
	}
}

func (p *Player) SendLeaderboard(s *Server) {

	leaderboard := Leaderboard{
		Current: p,
		Top:     []*Player{},
	}
	for _, player := range s.Players {
		if player.Iinitialized {
			rounded_player := *player
			rounded_player.Size = math.Round(rounded_player.Size)
			leaderboard.Top = append(leaderboard.Top, &rounded_player)
		}
	}

	rounded_current := *p
	rounded_current.Size = math.Round(rounded_current.Size)
	rounded_current.Pos.X = math.Round(rounded_current.Pos.X)
	rounded_current.Pos.Y = math.Round(rounded_current.Pos.Y)

	leaderboard.Current = &rounded_current

	sort.Sort(BySize(leaderboard.Top))

	p.SendTemplate(s, "leaderboard.tmpl.html", leaderboard)
}

func (p *Player) SendTemplate(s *Server, template string, data any) error {
	var buf bytes.Buffer
	writer := &buf
	err := s.Templates.ExecuteTemplate(writer, template, data)
	if err != nil {
		return err
	}
	p.Lock()
	p.Conn.WriteMessage(
		websocket.TextMessage,
		buf.Bytes())
	p.Unlock()
	return nil
}

func (p *Player) update(delta int64) {

	move_speed = 1.0 / (p.Size * 0.035)

	if p.Ctl.Up {
		p.Vel.Y -= move_speed * float64(delta)
	}
	if p.Ctl.Down {
		p.Vel.Y += move_speed * float64(delta)
	}
	if p.Ctl.Left {
		p.Vel.X -= move_speed * float64(delta)
	}
	if p.Ctl.Right {
		p.Vel.X += move_speed * float64(delta)
	}

	//Friction
	p.Vel.MultF(0.80)

	p.Vel.LimitF(30)

	p.Viewport.Pos.Add(p.Vel)

	p.Pos = p.Viewport.PlayerAbsolute()
	p.Pos.SubF(float64(p.Size / 2))

	p.Collider.X = p.Pos.X + ((p.Size * 0.50) / 2)
	p.Collider.Y = p.Pos.Y + ((p.Size * 0.50) / 2)

	p.Collider.Width = float64(p.Size * 0.50)
	p.Collider.Height = float64(p.Size * 0.50)

	p.EatPower = float64(p.Size) / 200

	if p.Size < 40 {
		p.Dead = true
	}

	if p.Size > 7000 {
		p.Dead = true
	}

	p.Size -= 0.1

	p.FontSize = p.Size * 0.25
}

func (p *Player) sendPlayer(s *Server) {
	p.SendTemplate(s, "self.tmpl.html", p)
}

func (p *Player) sendRenderer(s *Server) {
	p.SendTemplate(s, "renderer.tmpl.html", nil)
}
func (p *Player) sendDeadScreen(s *Server) {
	p.SendTemplate(s, "dead-screen.tmpl.html", nil)
}

func (p *Player) sendPostion(s *Server) {
	err := p.SendTemplate(s, "self-position.tmpl.html", p)
	if err != nil {
		log.Printf("Template error: %v", err)
	}
}
