package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	s := server.Server{
		Addr: ":8080",
		WsUpgrader: websocket.Upgrader{
			CheckOrigin:      func(r *http.Request) bool { return true },
			HandshakeTimeout: time.Second * 3,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		},
		GameHandler: server.NewLoggerHandler(logger, server.NewGameHandler(server.NewWaitList())),
	}

	log.Printf("server is running on %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
