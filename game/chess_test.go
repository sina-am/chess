package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)

	assert.ErrorIs(t,
		game.Play(
			"player1",
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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)

	game.Play(
		"player1",
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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)

	assert.Nil(t,
		game.Play(
			"player1",
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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)

	game.Print()
	assert.ErrorIs(t,
		game.Play(
			"player1",
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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)
	game.Print()

	t.Run("move to a place where king is checked", func(t *testing.T) {
		assert.ErrorIs(t,
			game.Play(
				"player1",
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
				"player1",
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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)
	game.Print()

	t.Run("move near king", func(t *testing.T) {
		assert.ErrorIs(t,
			game.Play(
				"player1",
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

	game := NewChessFromPieces(map[string]Color{"player1": White, "player2": Black}, pieces)
	game.Print()

	t.Run("a touch of death", func(t *testing.T) {
		assert.ErrorIs(t, game.Play(
			"player1",
			Move{
				From: Location{Row: 1, Col: 7},
				To:   Location{Row: 1, Col: 1},
			},
		), ErrGameEnd,
		)
		assert.Equal(t, game.GetWinner(), White)
	})
}
