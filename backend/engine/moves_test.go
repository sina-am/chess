package engine

// import (
// 	"testing"

// 	"github.com/sina-am/chess/types"
// 	"github.com/stretchr/testify/assert"
// )

// func TestIsValidPawnMove(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})

// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Pawn,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{5, 0},
// 			Location{6, 0},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Pawn,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{5, 0},
// 			Location{6, 1},
// 		),
// 	)
// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Pawn,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 0},
// 			Location{2, 1},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Pawn,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 0},
// 			Location{2, 0},
// 		),
// 	)
// }
// func TestIsValidQueenMove(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})

// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Queen,
// 				PieceNumber: 1,
// 				Color:       types.Black,
// 			},
// 			Location{7, 2},
// 			Location{2, 7},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Queen,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 4},
// 			Location{4, 1},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Queen,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 1},
// 			Location{4, 4},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Queen,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{4, 4},
// 			Location{1, 1},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Queen,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 1},
// 			Location{2, 5},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Queen,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 5},
// 			Location{2, 1},
// 		),
// 	)
// }
// func TestIsValidBishopMove(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})

// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Bishop,
// 				PieceNumber: 1,
// 				Color:       types.Black,
// 			},
// 			Location{7, 2},
// 			Location{2, 7},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Bishop,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 4},
// 			Location{4, 1},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Bishop,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 1},
// 			Location{4, 4},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Bishop,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{4, 4},
// 			Location{1, 1},
// 		),
// 	)
// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Bishop,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 1},
// 			Location{2, 5},
// 		),
// 	)
// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Bishop,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 5},
// 			Location{2, 1},
// 		),
// 	)
// }
// func TestIsValidKingMove(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})

// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.King,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 1},
// 			Location{2, 2},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.King,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 1},
// 			Location{2, 2},
// 		),
// 	)
// }
// func TestIsValidRookMove(t *testing.T) {
// 	game := NewStandardGame([2]string{"id1", "id2"})

// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Rook,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{1, 1},
// 			Location{2, 2},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Rook,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 1},
// 			Location{2, 2},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Rook,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 1},
// 			Location{2, 7},
// 		),
// 	)
// 	assert.True(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Rook,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 2},
// 			Location{6, 2},
// 		),
// 	)
// 	assert.False(t,
// 		game.isValidMove(
// 			&types.Piece{
// 				Type:        types.Rook,
// 				PieceNumber: 1,
// 				Color:       types.White,
// 			},
// 			Location{2, 2},
// 			Location{6, 4},
// 		),
// 	)
// }
