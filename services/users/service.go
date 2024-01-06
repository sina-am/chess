package users

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/core"
)

type Service interface {
	Start()
	SetMiddlewares(e *echo.Group)
	SetAPIs(e *echo.Echo)
}

type service struct {
	storage       Storage
	apis          *APIService
	authenticator Authenticator
}

func getStorage(cfg config.Config) (Storage, error) {
	if cfg.HasDatabase() {
		return NewMongoStorage(context.TODO(), cfg.GetDatabaseClient(), cfg.Database.Name)
	} else {
		return NewMemoryStorage(), nil
	}
}

func NewService(cfg config.Config) (*service, error) {
	storage, err := getStorage(cfg)
	if err != nil {
		return nil, err
	}
	authenticator := NewJWTAuthentication(cfg.SecretKey, storage)
	renderer, err := core.NewTemplateRenderer(cfg.Debug, "./services/users/templates")
	if err != nil {
		return nil, err
	}

	apis := NewAPIs(storage, authenticator, renderer)

	return &service{
		storage:       storage,
		apis:          apis,
		authenticator: authenticator,
	}, nil
}

func (s *service) Start() {
}

func (s *service) GetAuthenticator() Authenticator {
	return s.authenticator
}
func (s *service) SetAPIs(e *echo.Echo) {
	e.POST("/auth/login", s.apis.AuthenticationAPI)
	e.GET("/auth/login", s.apis.AuthenticationAPI)
	e.POST("/auth/registration", s.apis.RegistrationAPI)
	e.GET("/users", s.apis.UsersAPI)
}

func (s *service) SetMiddlewares(e *echo.Echo) {
}
