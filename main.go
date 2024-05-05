package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/services/auth"
	"github.com/sina-am/chess/services/game"
	"github.com/sina-am/chess/services/users"
	"github.com/sina-am/chess/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type App interface {
	AddRoutes(e *echo.Echo)
	SetRenderer(r core.Renderer)
	SetAuthenticator(a auth.Authenticator)
	GetName() string
	Start()
}

func startApps(cfg *config.Config, e *echo.Echo, apps []App, a auth.Authenticator) error {
	for _, app := range apps {
		r, err := core.NewTemplateRenderer(cfg.Debug, fmt.Sprintf("./services/%s/templates", app.GetName()))
		if err != nil {
			return err
		}
		app.SetRenderer(r)
		app.AddRoutes(e)
		app.SetAuthenticator(a)
		app.Start()
	}

	return nil
}

type userFetcher struct {
	storage storage.Storage
}

func (fetcher *userFetcher) GetUserById(ctx context.Context, id primitive.ObjectID) (auth.User, error) {
	return fetcher.storage.GetUserById(ctx, id)
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(1)
	e.Logger.SetHeader("${level}")

	// Middleware
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf_token",
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/static", "./static")

	e.Validator = core.NewValidator()

	cfg := &config.Config{
		Debug:           true,
		SecretKey:       "1234",
		DatabaseBackend: config.MongoBackend,
		Database: config.Database{
			Uri:      "mongodb://localhost",
			Username: "",
			Password: "",
			Name:     "chess",
			Timeout:  3 * time.Second,
		},
	}

	storage := storage.NewMemoryStorage()

	authenticator := auth.NewJWTAuthentication(cfg.SecretKey, &userFetcher{storage})
	userRenderer, err := core.NewTemplateRenderer(cfg.Debug, "./services/users/templates")
	if err != nil {
		log.Fatal(err)
	}
	userSrv := users.NewAPIService(storage, authenticator, userRenderer)

	e.POST("/auth/login", userSrv.AuthenticationPOST)
	e.GET("/auth/login", userSrv.AuthenticationGET)
	e.POST("/auth/registration", userSrv.RegistrationAPI)
	e.GET("/users", userSrv.UsersAPI)

	gameRenderer, err := core.NewTemplateRenderer(cfg.Debug, "./services/game/templates")
	if err != nil {
		log.Fatal(err)
	}
	gameSrv := game.NewAPIService(cfg, storage, authenticator, gameRenderer)

	e.GET("/players", gameSrv.GetPlayers)
	e.GET("/ws", gameSrv.WebSocketAPI)
	e.GET("/game-options", gameSrv.GameOptions)
	e.POST("/game-options", gameSrv.GameOptions)
	e.GET("/game", gameSrv.StartGame)
	e.GET("/", gameSrv.Home)

	go gameSrv.GameHandler.Start()

	e.Use(authenticator.AuthenticationMiddleware)
	e.Logger.Fatal(e.Start(":8080"))
}
