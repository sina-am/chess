package engine

// import (
// 	"testing"

// 	"github.com/sina-am/chess/types"
// 	"github.com/stretchr/testify/assert"
// )

// func TestChessNewGameBoard(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})
// 	game.Board.Print()
// }

// func TestChessChecked(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})
// 	game.Board[1][1] = nil
// 	game.Board[6][4] = game.Board[7][4]
// 	game.Board[7][4] = nil

// 	game.MovePiece(types.White, Location{0, 2}, Location{2, 0})
// 	assert.True(t, game.Turn.IsChecked)
// 	game.Board.Print()
// }
// func TestChessMove(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})
// 	// Not your turn
// 	err := game.MovePiece(types.Black, Location{1, 1}, Location{2, 2})
// 	assert.NotNil(t, err)

// 	// Invalid pawn move
// 	err = game.MovePiece(types.White, Location{1, 1}, Location{2, 2})
// 	assert.NotNil(t, err)

// 	// Valid pawn move
// 	err = game.MovePiece(types.White, Location{1, 1}, Location{2, 1})
// 	assert.Nil(t, err)
// 	assert.Equal(t, game.Board[2][1].Type, types.Pawn)
// 	assert.Nil(t, game.Board[1][1])

// 	// Valid Black knight move
// 	err = game.MovePiece(types.Black, Location{7, 1}, Location{5, 2})
// 	assert.Nil(t, err)
// 	assert.Equal(t, game.Board[5][2].Type, types.Knight)
// 	assert.Nil(t, game.Board[7][1])

// 	// Valid White Bishop move
// 	err = game.MovePiece(types.White, Location{0, 2}, Location{2, 0})
// 	assert.Nil(t, err)
// 	assert.Equal(t, game.Board[2][0].Type, types.Bishop)
// 	assert.Nil(t, game.Board[0][2])

// 	// Valid Black Rook move
// 	err = game.MovePiece(types.Black, Location{7, 0}, Location{7, 1})
// 	assert.Nil(t, err)
// 	assert.Equal(t, game.Board[7][1].Type, types.Rook)
// 	assert.Nil(t, game.Board[7][0])

// 	// Valid White Bishop take
// 	err = game.MovePiece(types.White, Location{2, 0}, Location{6, 4})
// 	assert.Nil(t, err)
// 	assert.Equal(t, game.Board[6][4].Type, types.Bishop)
// 	assert.Nil(t, game.Board[2][0])

// 	// Valid Black King take
// 	err = game.MovePiece(types.Black, Location{7, 4}, Location{6, 4})
// 	assert.Nil(t, err)
// 	game.Board.Print()
// }
