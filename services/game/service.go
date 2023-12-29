package game

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/services/users"
	"github.com/sina-am/chess/utils"
)

type service struct {
	apis *APIService
}

type Config struct {
	Authenticator users.Authenticator
	Debug         bool
}

func NewService(config Config) (*service, error) {
	renderer, err := utils.NewTemplateRenderer(config.Debug, "./services/game/templates")
	if err != nil {
		return nil, err
	}
	return &service{
		apis: &APIService{
			WsUpgrader: websocket.Upgrader{
				CheckOrigin:      func(r *http.Request) bool { return true },
				HandshakeTimeout: time.Second * 3,
				ReadBufferSize:   1024,
				WriteBufferSize:  1024,
			},
			GameHandler:   NewGameHandler(NewWaitList()),
			Authenticator: config.Authenticator,
			Renderer:      renderer,
		},
	}, nil
}

func (s *service) Start() {
	go s.apis.GameHandler.Start()
}

func (s *service) SetMiddlewares(e *echo.Echo) {
}
func (s *service) SetAPIs(e *echo.Echo) {
	e.GET("/ws", s.apis.WebSocketAPI)
	e.GET("/game-options", s.apis.GameOptions)
	e.POST("/game-options", s.apis.GameOptions)
	e.GET("/game", s.apis.StartGame)
	e.GET("/", s.apis.Home)
}
