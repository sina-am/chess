package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/apis"
	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/service"
	"github.com/sina-am/chess/types"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	config, err := NewConfig()
	if err != nil {
		logger.Fatal(err.Error())
	}

	db, err := database.NewMongoDatabase(context.TODO(), config.DBAddr, config.DBName)
	if err != nil {
		logger.Fatal(err.Error())
	}

	gameSrv, err := service.NewGameService(db)
	if err != nil {
		logger.Fatal(err.Error())
	}
	types.NewValidator()

	server := apis.APIServer{
		Addr:          config.SrvAddr,
		Game:          gameSrv,
		Logger:        logger.Sugar(),
		Database:      db,
		Authenticator: apis.NewJWTAuthentication(config.SecretKey, db),
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	logger.Fatal(server.Run().Error())
}
