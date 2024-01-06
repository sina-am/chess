package game

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/services/users"
)

type APIService struct {
	WsUpgrader    websocket.Upgrader
	GameHandler   GameHandler
	Renderer      core.Renderer
	Authenticator users.Authenticator
}

func (s *APIService) GameOptions(c echo.Context) error {
	return s.Renderer.Render(c, "game-options.html", nil)
}

type gameOptionsIn struct {
	Mode     string        `query:"game_mode" validate:"required,eq=online|eq=offline"`
	Duration time.Duration `query:"duration" validate:"required,gte=10"`
}

func (s *APIService) StartGame(c echo.Context) error {
	opts := gameOptionsIn{}
	if err := c.Bind(&opts); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	if err := c.Validate(&opts); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return s.Renderer.Render(c, "game.html", map[string]any{
		"gameOpts": opts,
		"user":     s.Authenticator.GetUser(c),
	})
}

func (s *APIService) Home(c echo.Context) error {
	user := s.Authenticator.GetUser(c)
	content := map[string]any{
		"user": user,
	}

	return s.Renderer.Render(c, "home.html", content)
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
