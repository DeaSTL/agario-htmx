package main

import (
	"jmhart.dev/htmx-argio/server"
	"log"
	"net/http"
)

var Server server.Server

type indexData struct {
	WsUrl string
}

func index(w http.ResponseWriter, r *http.Request) {
	Server.Templates.ExecuteTemplate(w, "index.tmpl.html", nil)
}

func renderer(w http.ResponseWriter, r *http.Request) {
	Server.Templates.ExecuteTemplate(w, "renderer.tmpl.html", nil)
}

func main() {

	Server = server.Server{}
	http.HandleFunc("/renderer", renderer)
	Server.New("/ws")
	http.HandleFunc("/", index)

	log.Println("Starting agario server")
	http.ListenAndServe(":8080", nil)
}
