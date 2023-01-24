package main

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/apis"
	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/engine"
	"github.com/sina-am/chess/types"
	"go.uber.org/zap"
)

func main() {
	config, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewMongoDatabase(context.TODO(), config.DBAddr, config.DBName)
	if err != nil {
		log.Fatal(err)
	}

	game, err := engine.NewStandardGame(db)
	if err != nil {
		log.Fatal(err)
	}
	types.NewValidator()

	server := apis.APIServer{
		Addr:          config.SrvAddr,
		Game:          game,
		Logger:        logger.Sugar(),
		Database:      db,
		Authenticator: apis.NewJWTAuthentication(config.SecretKey, db),
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
	logger.Fatal(server.Run().Error())
}
