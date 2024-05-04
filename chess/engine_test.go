package chess

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPawnMoves(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 7},
		},
		{
			Type:     Pawn,
			Color:    Black,
			Location: Location{Row: 6, Col: 6},
		},
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Pawn,
			Color:    White,
			Location: Location{Row: 1, Col: 1},
		},
	}
	game := NewFromPieces(pieces)
	game.Print()
	possibleMoves := game.possibleMoves[pieces[3]]

	for _, move := range possibleMoves {
		fmt.Printf("%d, %d\n", move.Row, move.Col)
	}
}
func TestBishopMoves(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 7},
		},
		{
			Type:     Queen,
			Color:    Black,
			Location: Location{Row: 6, Col: 6},
		},
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Bishop,
			Color:    White,
			Location: Location{Row: 4, Col: 4},
		},
	}
	game := NewFromPieces(pieces)
	game.Print()
	possibleMoves := game.possibleMoves[pieces[3]]

	for _, move := range possibleMoves {
		fmt.Printf("%d, %d\n", move.Row, move.Col)
	}
}

func TestKingMoves(t *testing.T) {
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
	}
	game := NewFromPieces(pieces)
	err := game.Play(White, Move{From: Location{Row: 0, Col: 0}, To: Location{Row: 1, Col: 1}})
	assert.Nil(t, err)
}
func TestKingProblemMoves(t *testing.T) {
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
			Color:    Black,
			Location: Location{Row: 1, Col: 6},
		},
	}
	game := NewFromPieces(pieces)
	err := game.Play(White, Move{From: Location{Row: 0, Col: 0}, To: Location{Row: 1, Col: 1}})
	assert.NotNil(t, err)
	err = game.Play(White, Move{From: Location{Row: 0, Col: 0}, To: Location{Row: 1, Col: 0}})
	assert.NotNil(t, err)
}
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
		ErrInvalidPieceMove,
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
			ErrInvalidPieceMove,
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
			ErrInvalidPieceMove,
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
		assert.Equal(t, game.GetResult().WinnerColor, White)
	})
}

func TestKingCapture(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 7, Col: 7},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Queen,
			Color:    White,
			Location: Location{Row: 1, Col: 1},
		},
	}
	game := NewFromPieces(pieces)
	game.SwitchTurn()

	err := game.Play(Black, Move{From: Location{Row: 0, Col: 0}, To: Location{Row: 1, Col: 1}})
	assert.Nil(t, err)
	game.Print()

}
func TestWhiteRightCastling(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 4},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 4},
		},
		{
			Type:     Rook,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Rook,
			Color:    White,
			Location: Location{Row: 0, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	assert.Nil(t, game.Play(White, Move{From: Location{Row: 0, Col: 4}, To: Location{Row: 0, Col: 6}}))

	assert.Equal(t, game.board[0][6], pieces[0])
	assert.Equal(t, game.board[0][5], pieces[3])
	game.Print()
}
func TestWhiteLeftCastling(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 4},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 4},
		},
		{
			Type:     Rook,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     Rook,
			Color:    White,
			Location: Location{Row: 0, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	assert.Nil(t, game.Play(White, Move{From: Location{Row: 0, Col: 4}, To: Location{Row: 0, Col: 2}}))
	assert.Equal(t, game.board[0][2], pieces[0])
	assert.Equal(t, game.board[0][3], pieces[2])
	game.Print()
}
func TestBlackRightCastling(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 4},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 4},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	game.turn = Black
	assert.Nil(t, game.Play(Black, Move{From: Location{Row: 7, Col: 4}, To: Location{Row: 7, Col: 6}}))
	assert.Equal(t, game.board[7][6], pieces[1])
	assert.Equal(t, game.board[7][5], pieces[3])
	game.Print()
}
func TestBlackLeftCastling(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 4},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 4},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	game.switchTurn()
	assert.Nil(t, game.Play(Black, Move{From: Location{Row: 7, Col: 4}, To: Location{Row: 7, Col: 2}}))

	assert.Equal(t, game.board[7][2], pieces[1])
	assert.Equal(t, game.board[7][3], pieces[2])
	game.Print()
}
func TestLeftCastlingRollback(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 4},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 4},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	game.turn = Black
	rb := NewRollBack(game)

	rb.Do(Move{From: Location{Row: 7, Col: 4}, To: Location{Row: 7, Col: 2}})
	rb.RollBack()
	assert.Equal(t, game.board[7][4], pieces[1])
	assert.Equal(t, game.board[7][0], pieces[2])
	game.Print()
}
func TestRightCastlingRollback(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 4},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 7, Col: 4},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 0},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	game.turn = Black
	rb := NewRollBack(game)

	rb.Do(Move{From: Location{Row: 7, Col: 4}, To: Location{Row: 7, Col: 6}})
	rb.RollBack()
	assert.Equal(t, game.board[7][4], pieces[1])
	assert.Equal(t, game.board[7][7], pieces[3])
	game.Print()
}
func TestStaleMate(t *testing.T) {
	pieces := []*Piece{
		{
			Type:     King,
			Color:    White,
			Location: Location{Row: 0, Col: 0},
		},
		{
			Type:     King,
			Color:    Black,
			Location: Location{Row: 2, Col: 2},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 7, Col: 1},
		},
		{
			Type:     Rook,
			Color:    Black,
			Location: Location{Row: 2, Col: 7},
		},
	}
	game := NewFromPieces(pieces)
	game.turn = Black

	game.Play(Black, Move{
		From: Location{
			Row: 2,
			Col: 7,
		},
		To: Location{
			Row: 1,
			Col: 7,
		},
	})

	assert.Equal(t, Stalemate, game.GetResult().Reason)
}
