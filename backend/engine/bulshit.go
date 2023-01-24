package engine

import (
	"github.com/sina-am/chess/types"
)

func makePieces() []types.Piece {
	return []types.Piece{
		{
			Type:        types.Rook,
			PieceNumber: 1,
			Color:       types.White,
		},
		{
			Type:        types.Knight,
			PieceNumber: 1,
			Color:       types.White,
		},
		{
			Type:        types.Bishop,
			PieceNumber: 1,
			Color:       types.White,
		},
		{
			Type:        types.Queen,
			PieceNumber: 1,
			Color:       types.White,
		},
		{
			Type:        types.King,
			PieceNumber: 1,
			Color:       types.White,
		},
		{
			Type:        types.Bishop,
			PieceNumber: 2,
			Color:       types.White,
		},
		{
			Type:        types.Knight,
			PieceNumber: 2,
			Color:       types.White,
		},
		{
			Type:        types.Rook,
			PieceNumber: 2,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 1,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 2,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 3,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 4,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 5,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 6,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 7,
			Color:       types.White,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 8,
			Color:       types.White,
		},

		// Black pieces
		{
			Type:        types.Rook,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.Knight,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.Bishop,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.Queen,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.King,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.Bishop,
			PieceNumber: 2,
			Color:       types.Black,
		},
		{
			Type:        types.Knight,
			PieceNumber: 2,
			Color:       types.Black,
		},
		{
			Type:        types.Rook,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 1,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 2,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 3,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 4,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 5,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 6,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 7,
			Color:       types.Black,
		},
		{
			Type:        types.Pawn,
			PieceNumber: 8,
			Color:       types.Black,
		},
	}
}
func makeBoard(pieces []types.Piece) types.Board {
	return types.Board{
		{&pieces[0], &pieces[1], &pieces[2], &pieces[3], &pieces[4], &pieces[5], &pieces[6], &pieces[7]},
		{&pieces[8], &pieces[9], &pieces[10], &pieces[11], &pieces[12], &pieces[13], &pieces[14], &pieces[15]},
		{nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil},
		{&pieces[24], &pieces[25], &pieces[26], &pieces[27], &pieces[28], &pieces[29], &pieces[30], &pieces[31]},
		{&pieces[16], &pieces[17], &pieces[18], &pieces[19], &pieces[20], &pieces[21], &pieces[22], &pieces[23]},
	}
}
