package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	Addr        string
	WsUpgrader  websocket.Upgrader
	GameHandler *gameHandler
}

func (s *Server) ListenAndServe() error {
	go s.GameHandler.Start()

	dir := http.Dir("./static")
	fs := http.FileServer(dir)

	mux := http.NewServeMux()
	mux.Handle("/", fs)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := s.WsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrader error: %v", err)
			return
		}
		log.Printf("new player %s connected", conn.RemoteAddr())

		p := NewPlayer(conn, s.GameHandler)

		p.gameHandler.register <- p

		go p.ReadConn()
		go p.WriteConn()
	})

	return http.ListenAndServe(s.Addr, mux)
}
