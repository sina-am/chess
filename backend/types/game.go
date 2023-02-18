package types

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GameType string

const (
	RandomGame   GameType = "RANDOM"
	FriendlyGame GameType = "FRIENDLY"
)

type PlayRequest struct {
	From Location `json:"from" validate:"required"`
	To   Location `json:"to" validate:"required"`
}

func (r *PlayRequest) Validate() error {
	if err := validate.Struct(r); err != nil {
		return err
	}

	if r.From.Col > 7 || r.From.Row > 7 || r.To.Col > 7 || r.To.Row > 7 {
		return fmt.Errorf("out of board")
	}

	return nil
}

type Player struct {
	UserId       primitive.ObjectID `json:"user_id" bson:"user_id"`
	Color        Color              `json:"color" bson:"color"`
	IsChecked    bool               `json:"is_checked" bson:"is_checked"`
	IsCheckmated bool               `json:"is_checkmated" bson:"is_checkmated"`
	Turn         bool               `json:"turn" bson:"turn"`
}

type Game struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	StartedAt  time.Time          `json:"started_at" bson:"started_at"`
	Duration   time.Duration      `json:"duration" bson:"duration"`
	Players    []*Player          `json:"players" bson:"players"`
	StartedBy  primitive.ObjectID `json:"started_by" bson:"started_by"` // User that start the game
	IsAccepted bool               `json:"is_accepted" bson:"is_accepted"`
	Pieces     []*Piece           `json:"pieces" bson:"pieces,omitempty"`
}

type StartGameRequest struct {
	Duration time.Duration `json:"duration"`
	Type     GameType      `json:"type"`
}

func (r *StartGameRequest) Validate() error {
	return validate.Struct(r)
}
