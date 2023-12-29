package game

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/services/users"
	"github.com/sina-am/chess/utils"
)

type APIService struct {
	WsUpgrader    websocket.Upgrader
	GameHandler   GameHandler
	Renderer      utils.Renderer
	Authenticator users.Authenticator
}

func (s *APIService) GameOptions(c echo.Context) error {
	return s.Renderer.Render(c.Response().Writer, "game-options.html", nil)
}

func (s *APIService) StartGame(c echo.Context) error {
	gameMode := c.QueryParam("gameMode")
	return s.Renderer.Render(c.Response().Writer, "game.html", map[string]string{
		"gameMode":   gameMode,
		"playerName": "Sina",
	})
}

func (s *APIService) Home(c echo.Context) error {
	return s.Renderer.Render(c.Response().Writer, "home.html", nil)
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
