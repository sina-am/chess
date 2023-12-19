package main

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sina-am/chess/services/game"
	"github.com/sina-am/chess/services/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost").SetTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("database error: %s", err.Error())
	}
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/", "./static")

	userService, err := users.NewService(users.Config{
		MongoClient: client,
		DBName:      "users",
		SecretKey:   "1234",
	})
	if err != nil {
		log.Fatal(err)
	}

	userService.Start()
	userService.SetAPIs(e)

	// g := e.Group("/users")
	// userService.SetMiddlewares(g)

	gameService, err := game.NewService(game.Config{
		Authenticator: userService.GetAuthenticator(),
	})
	if err != nil {
		log.Fatal(err)
	}

	gameService.Start()
	gameService.SetAPIs(e)

	e.Logger.Fatal(e.Start(":8080"))
}
