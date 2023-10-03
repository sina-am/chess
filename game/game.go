package game

import (
	"errors"
)

var ErrGameEnd = errors.New("game ended")

type Move struct {
	From Location `json:"from"`
	To   Location `json:"to"`
}

func (m Move) Validate() error {
	return nil
}

type Chess interface {
	GetWinner() (string, error)
	InGame(playerId string) bool
	GetPlayers() []string
	Play(playerId string, m Move) error
	GetPlayerColor(playerId string) Color
	Exit(playerId string) (string, error)
}
