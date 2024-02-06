package users

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/auth"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type app struct {
	apis *APIService
}
type userFetcher struct {
	storage storage.Storage
}

func (fetcher *userFetcher) GetUserById(ctx context.Context, id primitive.ObjectID) (auth.User, error) {
	return fetcher.storage.GetUserById(ctx, id)
}

func NewApp(cfg *config.Config, storage storage.Storage) (*app, error) {
	authenticator := auth.NewJWTAuthentication(cfg.SecretKey, &userFetcher{storage})

	apis := &APIService{
		Authenticator: authenticator,
		Storage:       storage,
	}

	return &app{
		apis: apis,
	}, nil
}

func (s *app) Start() {
}

func (s *app) GetName() string {
	return "users"
}

func (s *app) SetRenderer(r core.Renderer) {
	s.apis.Renderer = r
}

func (s *app) SetAuthenticator(auth.Authenticator) {}
func (s *app) GetAuthenticator() auth.Authenticator {
	return s.apis.Authenticator
}

func (s *app) AddRoutes(e *echo.Echo) {
	e.POST("/auth/login", s.apis.AuthenticationAPI)
	e.GET("/auth/login", s.apis.AuthenticationAPI)
	e.POST("/auth/registration", s.apis.RegistrationAPI)
	e.GET("/users", s.apis.UsersAPI)
}
