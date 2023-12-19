package chess

import (
	"errors"
)

var (
	ErrGameEnd          = errors.New("game ended")
	ErrOutOfBoardMove   = errors.New("move should be between [0, 8)")
	ErrChecked          = errors.New("error checked cant move there")
	ErrInvalidPieceMove = errors.New("piece can't move like that")
	ErrNotPlayersTurn   = errors.New("it's not your turn")
)

type Move struct {
	From Location `json:"from"`
	To   Location `json:"to"`
}

func (m Move) Validate() error {
	if err := m.From.Validate(); err != nil {
		return err
	}
	if err := m.To.Validate(); err != nil {
		return err
	}
	return nil
}

type Chess interface {
	GetWinner() Color
	IsFinished() bool
	Play(playerId string, m Move) error
	Exit(playerId string) (Color, error)

	InGame(playerId string) bool
	GetPlayers() []string
	GetPlayerByColor(c Color) string
}
