package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPromotion(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 5},
		},
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Pawn,
			Color:    White,
			Location: Location{Row: 6, Col: 1},
		},
	}

	game := NewFromPieces(pieces)
	t.Run("promote pawn", func(t *testing.T) {
		assert.Nil(t, game.Play(White, Move{From: Location{Row: 6, Col: 1}, To: Location{Row: 7, Col: 1}}))
	})
	t.Run("king checked", func(t *testing.T) {
		assert.ErrorIs(t, game.Play(Black, Move{From: Location{Row: 7, Col: 5}, To: Location{Row: 7, Col: 4}}), ErrChecked)
	})
	game.Print()
}

func TestRookMove(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 5},
		},
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Rook,
			Color:    White,
			Location: Location{Row: 4, Col: 4},
		},
	}

	game := NewFromPieces(pieces)
	game.Print()

	t.Run("move rook to left", func(t *testing.T) {
		assert.Nil(t, game.Play(White, Move{From: Location{Row: 4, Col: 4}, To: Location{Row: 4, Col: 3}}))
	})

	game.switchTurn()
	t.Run("move rook to right", func(t *testing.T) {
		assert.Nil(t, game.Play(White, Move{From: Location{Row: 4, Col: 3}, To: Location{Row: 4, Col: 7}}))
	})
	game.switchTurn()
	t.Run("move rook up", func(t *testing.T) {
		assert.Nil(t, game.Play(White, Move{From: Location{Row: 4, Col: 7}, To: Location{Row: 1, Col: 7}}))
	})
	game.switchTurn()
	t.Run("move rook down", func(t *testing.T) {
		assert.Nil(t, game.Play(White, Move{From: Location{Row: 1, Col: 7}, To: Location{Row: 7, Col: 7}}))
	})
}
func TestKnightMove(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 5},
		},
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Knight,
			Color:    White,
			Location: Location{Row: 2, Col: 1},
		},
	}

	game := NewFromPieces(pieces)
	game.Print()

	err := game.Play(White, Move{From: Location{Row: 2, Col: 1}, To: Location{Row: 0, Col: 0}})
	assert.Error(t, err)
}
func TestCaptureOwnPiece(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 5},
		},
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Queen,
			Color:    White,
			Location: Location{Row: 1, Col: 0},
		},
	}

	game := NewFromPieces(pieces)

	assert.ErrorIs(t,
		game.Play(
			White,
			Move{
				From: Location{Row: 0, Col: 0},
				To:   Location{Row: 1, Col: 0},
			},
		), ErrInvalidPieceMove,
	)
}
func TestCapturePiece(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 2},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 3, Col: 3},
		},
		{
			Type:     Queen,
			Color:    White,
			Location: Location{Row: 6, Col: 0},
		},
		{
			Type:     Queen,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
	}

	game := NewFromPieces(pieces)

	game.Play(
		White,
		Move{
			From: Location{Row: 6, Col: 0},
			To:   Location{Row: 7, Col: 0},
		},
	)

	assert.True(t, pieces[3].Captured)
}
func TestSomething(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 3, Col: 3},
		},
		{
			Type:     Queen,
			Color:    White,
			Location: Location{Row: 6, Col: 0},
		},
		{
			Type:     Queen,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
	}

	game := NewFromPieces(pieces)

	assert.Nil(t,
		game.Play(
			White,
			Move{
				From: Location{Row: 6, Col: 0},
				To:   Location{Row: 7, Col: 0},
			},
		),
	)

	game.Print()

}

func TestPinnedPiece(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 3, Col: 3},
		},
		{
			Type:     Queen,
			Color:    White,
			Location: Location{Row: 1, Col: 0},
		},
		{
			Type:     Queen,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
	}

	game := NewFromPieces(pieces)

	game.Print()
	assert.ErrorIs(t,
		game.Play(
			White,
			Move{
				From: Location{Row: 1, Col: 0},
				To:   Location{Row: 1, Col: 1},
			},
		),
		ErrChecked,
	)
}

func TestKingMoveWhenIsChecked(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 3, Col: 3},
		},
		{
			Type:     Queen,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
	}

	game := NewFromPieces(pieces)
	game.Print()

	t.Run("move to a place where king is checked", func(t *testing.T) {
		assert.ErrorIs(t,
			game.Play(
				White,
				Move{
					From: Location{Row: 0, Col: 0},
					To:   Location{Row: 1, Col: 0},
				},
			),
			ErrChecked,
		)
	})
	t.Run("move to a safe squire", func(t *testing.T) {
		assert.Nil(t,
			game.Play(
				White,
				Move{
					From: Location{Row: 0, Col: 0},
					To:   Location{Row: 1, Col: 1},
				},
			),
		)
	})
}

func TestKingCollisions(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 3, Col: 1},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 3, Col: 3},
		},
	}

	game := NewFromPieces(pieces)
	game.Print()

	t.Run("move near king", func(t *testing.T) {
		assert.ErrorIs(t,
			game.Play(
				White,
				Move{
					From: Location{Row: 3, Col: 1},
					To:   Location{Row: 3, Col: 2},
				},
			),
			ErrChecked,
		)
	})
}
func TestKingCheckmate(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 2, Col: 0},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Queen,
			Color:    White,
			Location: Location{Row: 1, Col: 7},
		},
	}

	game := NewFromPieces(pieces)
	game.Print()

	t.Run("a touch of death", func(t *testing.T) {
		game.Play(
			White,
			Move{
				From: Location{Row: 1, Col: 7},
				To:   Location{Row: 1, Col: 1},
			},
		)
		assert.Equal(t, game.GetWinner(), White)
	})
}
