package game

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type APIService struct {
	WsUpgrader  websocket.Upgrader
	GameHandler GameHandler
}

func (s *APIService) WebSocketAPI(c echo.Context) error {
	conn, err := s.WsUpgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return err
	}
	log.Printf("new player %s connected", conn.RemoteAddr())

	p := NewPlayer(conn, s.GameHandler)

	p.gameHandler.Register(p)

	go p.ReadConn()
	go p.WriteConn()
	return nil
}
