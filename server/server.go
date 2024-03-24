package server

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"jmhart.dev/htmx-argio/physics"
	"jmhart.dev/htmx-argio/utils"
)

const GridWidth int = 10000
const GridHeight int = 10000

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Request struct {
	Request     string `json:"HX-Request"`
	Trigger     string `json:"HX-Trigger"`
	TriggerName string `json:"HX-Trigger-Name"`
	Target      string `json:"HX-Target"`
	CurrentURL  string `json:"HX-Current-URL"`
}

type Screen struct {
	Width  int
	Height int
}

type RawRequest struct {
	Headers Request `json:"HEADERS"`
	Screen  Screen  `json:"screen"`
}

type Server struct {
	Players   map[string]*Player
	Food      map[string]*Food
	Templates *template.Template
}

func (s *Server) initWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	log.Printf("Iniitalizing web socket")
	s.NewConnection(conn)
}

func (s *Server) New(path string) {
	http.HandleFunc("/ws", s.initWebSocket)
	s.Players = map[string]*Player{}
	s.Food = map[string]*Food{}
	s.Templates = template.New("")
	s.Templates.ParseGlob("templates/*")

	for i := 0; i < 1000; i++ {
		id := "food_" + utils.GenID(64)
		new_food := Food{
			ID:    id,
			Color: utils.GenerateRandomHexColor(),
			Pos:   physics.Vec2f{X: rand.Float64() * float64(GridWidth), Y: rand.Float64() * float64(GridHeight)},
		}
		new_food.New()
		s.Food[id] = &new_food
	}

	go func(s *Server) {
		var delta int64
		for {
			startTime := time.Now()
			for _, player := range s.Players {
				player.update(delta)
				player.sendPostion(s)
				for _, food := range s.Food {
					if player.Collider.IsColliding(&food.Collider) && !food.Consumed {
						food.Consume(player, s)
					}
				}
			}
			s.sendPlayerPositions(nil)
			time.Sleep(time.Millisecond * 150)
			endTime := time.Now()
			delta = endTime.Sub(startTime).Milliseconds()
		}
	}(s)
}

func (s *Server) updatePlayerGlobs(player *Player) {
	var buf bytes.Buffer
	writer := &buf

	s.Templates.ExecuteTemplate(writer, "players.tmpl.html", s.Players)
	s.Broadcast(
		websocket.TextMessage,
		string(buf.Bytes()))
}

func (s *Server) sendFood() {
	s.BroadcastTemplate("food.tmpl.html", s.Food)
}

func (s *Server) sendPlayerPositions(player *Player) {
	s.BroadcastTemplate("player-positions.tmpl.html", s.Players)
}

func (s *Server) handleMessage(req *RawRequest, player *Player) {
	switch req.Headers.Target {
	// Key Ups
	case "key-up-up":
		player.Ctl.Up = true
	case "key-down-up":
		player.Ctl.Down = true
	case "key-left-up":
		player.Ctl.Left = true
	case "key-right-up":
		player.Ctl.Right = true
		// Key downs
	case "key-up-down":
		player.Ctl.Up = false
	case "key-down-down":
		player.Ctl.Down = false
	case "key-left-down":
		player.Ctl.Left = false
	case "key-right-down":
		player.Ctl.Right = false
	case "viewport-resize":
		player.Viewport.Width = req.Screen.Width
		player.Viewport.Height = req.Screen.Height
		s.sendPlayerPositions(nil)
	}
}

func (s *Server) NewConnection(conn *websocket.Conn) {
	id := "player_" + utils.GenID(16)
	new_player := Player{}
	new_player.New(conn, id)
	new_player.Pos.X = (rand.Float64() - 0.5) * 400
	new_player.Pos.Y = (rand.Float64() - 0.5) * 400
	s.Players[id] = &new_player
	conn.SetCloseHandler(s.LostConnectionHandler(&new_player))

	s.updatePlayerGlobs(&new_player)
	s.sendPlayerPositions(&new_player)
	s.sendFood()
	new_player.sendPlayer(s)
	new_player.sendPostion(s)

	go func(player *Player) {

		for {
			msgType, msg, err := player.Conn.ReadMessage()
			if err != nil {
				return
			}

			if msgType == websocket.TextMessage {

				raw_request := RawRequest{}

				json.Unmarshal(msg, &raw_request)

				s.handleMessage(&raw_request, player)
			}
		}
	}(&new_player)
}

func (s *Server) LostConnectionHandler(player *Player) func(code int, text string) error {
	return func(code int, text string) error {
		delete(s.Players, player.ID)
		s.updatePlayerGlobs(nil)
		s.sendPlayerPositions(nil)
		return nil
	}
}

func (s *Server) BroadcastTemplate(template string, data any) {
	for _, player := range s.Players {
		err := player.SendTemplate(s, template, data)
		if err != nil {
			log.Printf("Tempalte Error: %v", err)
		}
	}
}

func (s Server) Broadcast(messageType int, message string) {
	for _, player := range s.Players {
		err := player.Conn.WriteMessage(messageType, []byte(message))
		if err != nil {
			log.Printf("Error broadcasting %+v", err)
		}
	}
}
