package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

// Upgrader to handle WebSocket connection upgrades
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin
		return true
	},
}

// handleWS handles incoming WebSocket connections
func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	fmt.Println("New incoming connection:", ws.RemoteAddr())

	s.conns[ws] = true

	// Fetch data from API
	resp := apiRequest()

	// Send response to WebSocket client
	err = ws.WriteMessage(websocket.TextMessage, []byte(resp))
	if err != nil {
		log.Printf("Failed to write message to WebSocket: %v", err)
	}
}

func apiRequest() string {
	resp, err := http.Get(API_KEY_HERE)
	if err != nil {
		log.Printf("error fetching data from Xano: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response body: %v", err)
	}

	return string(body)
}

func main() {
	server := NewServer()
	http.HandleFunc("/ws", server.handleWS)
	log.Printf("Server started on :55555")
	if err := http.ListenAndServe(":55555", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
