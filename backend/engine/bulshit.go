package engine

import (
	"github.com/sina-am/chess/types"
)

func makePieces() []*types.Piece {
	return []*types.Piece{
		{
			Type:        types.Rook,
			PieceNumber: 1,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 0},
		},
		{
			Type:        types.Knight,
			PieceNumber: 1,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 1},
		},
		{
			Type:        types.Bishop,
			PieceNumber: 1,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 2},
		},
		{
			Type:        types.Queen,
			PieceNumber: 1,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 3},
		},
		{
			Type:        types.King,
			PieceNumber: 1,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 4},
		},
		{
			Type:        types.Bishop,
			PieceNumber: 2,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 5},
		},
		{
			Type:        types.Knight,
			PieceNumber: 2,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 6},
		},
		{
			Type:        types.Rook,
			PieceNumber: 2,
			Color:       types.White,
			Location:    types.Location{Row: 0, Col: 7},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 1,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 0},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 2,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 1},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 3,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 2},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 4,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 3},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 5,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 4},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 6,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 5},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 7,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 6},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 8,
			Color:       types.White,
			Location:    types.Location{Row: 1, Col: 7},
		},

		// Black pieces
		{
			Type:        types.Rook,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 0},
		},
		{
			Type:        types.Knight,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 1},
		},
		{
			Type:        types.Bishop,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 2},
		},
		{
			Type:        types.Queen,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 3},
		},
		{
			Type:        types.King,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 4},
		},
		{
			Type:        types.Bishop,
			PieceNumber: 2,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 5},
		},
		{
			Type:        types.Knight,
			PieceNumber: 2,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 6},
		},
		{
			Type:        types.Rook,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 7, Col: 7},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 1,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 0},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 2,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 1},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 3,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 2},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 4,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 3},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 5,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 4},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 6,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 5},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 7,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 6},
		},
		{
			Type:        types.Pawn,
			PieceNumber: 8,
			Color:       types.Black,
			Location:    types.Location{Row: 6, Col: 7},
		},
	}
}
