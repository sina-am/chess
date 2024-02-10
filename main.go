package main

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sina-am/chess/auth"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/services/game"
	"github.com/sina-am/chess/services/users"
	"github.com/sina-am/chess/storage"
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

	gameApp, err := game.NewApp(cfg, storage)
	if err != nil {
		e.Logger.Fatal(err)
	}
	userApp, err := users.NewApp(cfg, storage)
	if err != nil {
		e.Logger.Fatal(err)
	}
	authenticator := userApp.GetAuthenticator()
	startApps(cfg, e, []App{gameApp, userApp}, authenticator)

	e.Use(authenticator.AuthenticationMiddleware)
	e.Logger.Fatal(e.Start(":8080"))
}
