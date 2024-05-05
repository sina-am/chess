package game

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/services/auth"
	"github.com/sina-am/chess/storage"
)

type APIService struct {
	Storage       storage.Storage
	WsUpgrader    websocket.Upgrader
	GameHandler   GameHandler
	Renderer      core.Renderer
	Authenticator auth.Authenticator
}

func NewAPIService(cfg *config.Config, s storage.Storage, auth auth.Authenticator, renderer core.Renderer) *APIService {
	return &APIService{
		Storage: s,
		WsUpgrader: websocket.Upgrader{
			CheckOrigin:      func(r *http.Request) bool { return true },
			HandshakeTimeout: time.Second * 3,
			ReadBufferSize:   1024,
			WriteBufferSize:  1024,
		},
		GameHandler:   NewGameHandler(NewMemoryWaitList(), s),
		Authenticator: auth,
		Renderer:      renderer,
	}
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
func (s *APIService) GetPlayers(c echo.Context) error {
	users, err := s.Storage.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
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
	c.Logger().Debug("new player %s connected", conn.RemoteAddr())

	user := s.Authenticator.GetUser(c)
	p := NewWSClient(conn, s.GameHandler, user)

	p.gameHandler.Register(p, user)
	p.StartLoop(c.Request().Context())
	return nil
}
