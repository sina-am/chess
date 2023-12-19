package users

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
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

type Config struct {
	MongoClient *mongo.Client
	DBName      string
	SecretKey   string
}

func NewService(config Config) (*service, error) {
	storage, err := NewMongoStorage(context.TODO(), config.MongoClient, config.DBName)
	if err != nil {
		return nil, err
	}

	authenticator := NewJWTAuthentication(config.SecretKey, storage)
	apis := NewAPIs(storage, authenticator)

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
	e.POST("/auth/registration", s.apis.RegistrationAPI)
	e.GET("/users", s.apis.UsersAPI)
}

func (s *service) SetMiddlewares(e *echo.Group) {
	e.Use(s.apis.AuthenticationMiddleware)
}
