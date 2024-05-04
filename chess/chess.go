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
	GetResult() Result // Return game result if game is finished otherwise return NoResult
	Play(playerColor Color, m Move) error
	Exit() // Clear the game state
}
