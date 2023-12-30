package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sina-am/chess/config"
	"github.com/sina-am/chess/services/game"
	"github.com/sina-am/chess/services/users"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(1)
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/static", "./static")

	cfg := config.Config{
		Debug:     true,
		SecretKey: "1234",
		Database: config.Database{
			Uri:      "mongodb://localhost",
			Username: "",
			Password: "",
			Name:     "",
			Timeout:  3 * time.Second,
		},
	}

	userService, err := users.NewService(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userService.Start()
	userService.SetAPIs(e)
	userService.SetMiddlewares(e)

	gameService, err := game.NewService(cfg, userService.GetAuthenticator())
	if err != nil {
		log.Fatal(err)
	}

	gameService.Start()
	gameService.SetAPIs(e)

	e.Logger.Fatal(e.Start(":8080"))
}
