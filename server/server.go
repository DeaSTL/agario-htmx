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
	"jmhart.dev/htmx-argio/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Request struct {
	Raw         string
	Request     string `json:"HX-Request"`
	Trigger     string `json:"HX-Trigger"`
	TriggerName string `json:"HX-Trigger-Name"`
	Target      string `json:"HX-Target"`
	CurrentURL  string `json:"HX-Current-URL"`
}

type Server struct {
	Players   map[string]*Player
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
	s.Templates = template.New("")
	s.Templates.ParseGlob("templates/*")

	go func(s *Server) {
		for {
			for _, player := range s.Players {
				player.update()
				player.sendPlayerPostion(s)
			}
			s.sendPlayerPositions()
			time.Sleep(time.Millisecond * 300)
		}
	}(s)
}

func (s *Server) updatePlayerGlobs() {
	var buf bytes.Buffer
	writer := &buf
	s.Templates.ExecuteTemplate(writer, "players.tmpl.html", s.Players)
	s.Broadcast(
		websocket.TextMessage,
		string(buf.Bytes()))
}

func (s *Server) sendPlayerPositions() {
	var buf bytes.Buffer
	writer := &buf
	s.Templates.ExecuteTemplate(writer, "player-positions.tmpl.html", s.Players)
	s.Broadcast(
		websocket.TextMessage,
		string(buf.Bytes()))
}

func (s *Server) handleMessage(req *Request, player *Player) {
	switch req.Trigger {
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

	s.updatePlayerGlobs()
	s.sendPlayerPositions()
	new_player.sendPlayer(s)
	new_player.sendPlayerPostion(s)

	go func(player *Player) {

		for {
			msgType, msg, err := player.Conn.ReadMessage()
			if err != nil {
				return
			}

			if msgType == websocket.TextMessage {

				request_map := map[string]Request{}

				json.Unmarshal(msg, &request_map)

				request := request_map["HEADERS"]

				request.Raw = string(msg)

				s.handleMessage(&request, player)
			}
		}
	}(&new_player)
}

func (s *Server) LostConnectionHandler(player *Player) func(code int, text string) error {
	return func(code int, text string) error {
		delete(s.Players, player.ID)
		return nil
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
