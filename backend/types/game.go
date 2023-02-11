package types

import (
	"time"
)

type MovePieceRequest struct {
	From Location `json:"from"`
	To   Location `json:"to"`
}

type Player struct {
	UserId       string `json:"user_id" bson:"user_id"`
	Color        Color  `json:"color" bson:"color"`
	IsChecked    bool   `json:"is_checked" bson:"is_checked"`
	IsCheckmated bool   `json:"is_checkmated" bson:"is_checkmated"`
	Turn         bool   `json:"turn" bson:"turn"`
}

type Game struct {
	Id         string        `json:"_id" bson:"_id"`
	StartedAt  time.Time     `json:"started_at" bson:"started_at"`
	Duration   time.Duration `json:"duration" bson:"duration"`
	Players    []Player      `json:"players" bson:"players"`
	StartedBy  string        `json:"started_by" bson:"started_by"` // User that start the game
	IsAccepted bool          `json:"is_accepted" bson:"is_accepted"`
	Pieces     []*Piece      `json:"pieces" bson:"pieces,omitempty"`
}

type StartGameRequest struct {
	PlayerUserId string        `json:"player_user_id" validate:"required"`
	Duration     time.Duration `json:"duration"`
}

func (r *StartGameRequest) Validate() error {
	return validate.Struct(r)
}
