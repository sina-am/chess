package game

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/auth"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/storage"
)

type app struct {
	apis *APIService
}

func NewApp(cfg *config.Config, s storage.Storage) (*app, error) {
	return &app{
		apis: &APIService{
			WsUpgrader: websocket.Upgrader{
				CheckOrigin:      func(r *http.Request) bool { return true },
				HandshakeTimeout: time.Second * 3,
				ReadBufferSize:   1024,
				WriteBufferSize:  1024,
			},
			GameHandler:   NewGameHandler(NewMemoryWaitList()),
			Authenticator: nil,
			Renderer:      nil,
		},
	}, nil
}

func (s *app) GetName() string {
	return "game"
}

func (s *app) SetAuthenticator(auth auth.Authenticator) {
	s.apis.Authenticator = auth
}
func (s *app) SetRenderer(r core.Renderer) {
	s.apis.Renderer = r
}
func (s *app) Start() {
	go s.apis.GameHandler.Start()
}

func (s *app) AddRoutes(e *echo.Echo) {
	e.GET("/players", s.apis.GetPlayers)
	e.GET("/ws", s.apis.WebSocketAPI)
	e.GET("/game-options", s.apis.GameOptions)
	e.POST("/game-options", s.apis.GameOptions)
	e.GET("/game", s.apis.StartGame)
	e.GET("/", s.apis.Home)

}
