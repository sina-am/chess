package game

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/services/users"
)

type service struct {
	ws            *APIService
	authenticator users.Authenticator
}

type Config struct {
	Authenticator users.Authenticator
}

func NewService(config Config) (*service, error) {
	return &service{
		ws: &APIService{
			WsUpgrader: websocket.Upgrader{
				CheckOrigin:      func(r *http.Request) bool { return true },
				HandshakeTimeout: time.Second * 3,
				ReadBufferSize:   1024,
				WriteBufferSize:  1024,
			},
			GameHandler: NewGameHandler(NewWaitList()),
		},
		authenticator: config.Authenticator,
	}, nil
}

func (s *service) Start() {
	go s.ws.GameHandler.Start()
}

func (s *service) SetMiddlewares(e *echo.Echo) {
}
func (s *service) SetAPIs(e *echo.Echo) {
	e.GET("/ws", s.ws.WebSocketAPI)
}
