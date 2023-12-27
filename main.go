package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sina-am/chess/services/game"
	"github.com/sina-am/chess/services/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDatabase(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost").SetTimeout(3*time.Second))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("database error: %s", err.Error())
	}
	return client, nil
}
func setupWithoutDatabase(e *echo.Echo) {
	gameService, err := game.NewService(game.Config{
		Authenticator: nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	gameService.Start()
	gameService.SetAPIs(e)
	e.Logger.Warn("Running without database connection")
}

func setupWithDatabase(e *echo.Echo, client *mongo.Client) {
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
}

func main() {
	e := echo.New()
	e.Logger.SetLevel(1)
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Static("/", "./static")

	client, err := ConnectDatabase(context.Background())
	if err != nil {
		setupWithoutDatabase(e)
	} else {
		setupWithDatabase(e, client)
	}
	e.Logger.Fatal(e.Start(":8080"))
}
