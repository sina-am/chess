package types

import (
	"github.com/sina-am/chess/chess"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientEventType string

const (
	StartedClientEvent ClientEventType = "started"
	EndGameClientEvent ClientEventType = "ended"
	PlayedClientEvent  ClientEventType = "played"
)

type ServerEventType string

const (
	StartServerEvent        ServerEventType = "start"
	PlayServerEvent         ServerEventType = "play"
	ExitServerEvent         ServerEventType = "exit"
	OfferDrawServerEvent    ServerEventType = "offerDraw"
	ResponseDrawServerEvent ServerEventType = "respondDraw"
)

type StartGameMsgIn struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

type StartGameMsgOut struct {
	Type    ClientEventType        `json:"type"`
	Payload StartGamePayloadMsgOut `json:"payload"`
}

type StartGamePayloadMsgOut struct {
	Opponent Player `json:"opponent"`
	You      Player `json:"you"`
}

type PlayGameMsgIn struct {
	Move chess.Move `json:"move"`
}

type PlayGameMsgOut struct {
	Type    ClientEventType       `json:"type"`
	Payload PlayGamePayloadMsgOut `json:"payload"`
}

type PlayGamePayloadMsgOut struct {
	Player primitive.ObjectID `json:"player"`
	Move   chess.Move         `json:"move"`
}

type EndGameMsgOut struct {
	Type    ClientEventType      `json:"type"`
	Payload EndGamePayloadMsgOut `json:"payload"`
}

type EndGamePayloadMsgOut struct {
	Winner chess.Color  `json:"winner"`
	Score  int          `json:"score"`
	Reason chess.Reason `json:"reason"`
}
