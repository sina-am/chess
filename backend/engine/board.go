// Handles base movement of pieces and basic rules of chess standardBoard
package engine

import (
	"fmt"
	"math"

	"github.com/sina-am/chess/types"
)

type standardBoard struct {
	board [8][8]*types.Piece
	kings map[types.Color]*types.Piece
}

type Board interface {
	MovePiece(src, dst types.Location) error
	GetPiece(loc types.Location) (*types.Piece, error)
	GetAllPieces() []*types.Piece
	Print()
}

func NewStandardBoard() *standardBoard {
	return newStandardBoard()
}

func NewstandardBoardFromPieces(pieces []*types.Piece) *standardBoard {
	b := &standardBoard{
		kings: make(map[types.Color]*types.Piece, 2),
	}
	for _, piece := range pieces {
		if !piece.IsDead {
			b.board[piece.Location.Row][piece.Location.Col] = piece
			if piece.Type == types.King {
				b.kings[piece.Color] = piece
			}
		}
	}
	return b
}

func (b *standardBoard) Print() {
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.board[i][j] != nil {
				fmt.Printf("%s ", b.board[i][j].String())
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
}

func (b *standardBoard) GetAllPieces() []*types.Piece {
	pieces := make([]*types.Piece, 32)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if b.board[i][j] != nil {
				pieces = append(pieces, b.board[i][j])
			}
		}
	}
	return pieces
}

func (b *standardBoard) GetPiece(loc types.Location) (*types.Piece, error) {
	piece := b.board[loc.Row][loc.Col]
	if piece == nil {
		return nil, fmt.Errorf("invalid location")
	}
	return piece, nil
}

func (b *standardBoard) MovePiece(src, dst types.Location) error {
	if !b.isValidMove(src, dst) {
		return fmt.Errorf("invalid move")
	}
	if b.board[dst.Row][dst.Col] != nil {
		deadPiece := b.board[dst.Row][dst.Col]
		deadPiece.IsDead = true
	}

	b.board[dst.Row][dst.Col] = b.board[src.Row][src.Col]
	b.board[src.Row][src.Col] = nil

	b.board[dst.Row][dst.Col].Location.Col = dst.Col
	b.board[dst.Row][dst.Col].Location.Row = dst.Row

	return nil
}

func (b *standardBoard) isValidKingMove(piece *types.Piece, dst types.Location) bool {
	src := piece.Location
	return (src.Col == dst.Col &&
		math.Abs(float64(src.Row)-float64(dst.Row)) == 1) ||
		(src.Row == dst.Row &&
			math.Abs(float64(src.Col)-float64(dst.Col)) == 1)
}

func (b *standardBoard) isValidRookMove(piece *types.Piece, dst types.Location) bool {
	src := piece.Location
	if src.Col == dst.Col && src.Row < dst.Row {
		// Move up
		for rowStep := src.Row + 1; rowStep < dst.Row; rowStep++ {
			if b.board[rowStep][dst.Col] != nil {
				return false
			}
		}
		return true
	} else if src.Col == dst.Col && src.Row > dst.Row {
		// Move down
		for rowStep := src.Row - 1; rowStep > dst.Row; rowStep-- {
			if b.board[rowStep][dst.Col] != nil {
				return false
			}
		}
	} else if src.Row == dst.Row && src.Col < dst.Col {
		// Move right
		for colStep := src.Col + 1; colStep < dst.Col; colStep++ {
			if b.board[src.Row][colStep] != nil {
				return false
			}
		}
		return true
	} else if src.Row == dst.Row && src.Col > dst.Col {
		// Move left
		for colStep := src.Col - 1; colStep > dst.Col; colStep++ {
			if b.board[src.Col][colStep] != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (b *standardBoard) isValidBishopMove(piece *types.Piece, dst types.Location) bool {
	src := piece.Location
	if (dst.Col-src.Col) == (dst.Row-src.Row) && (dst.Col-src.Col) > 0 {
		// Move up-right
		colStep := src.Col + 1
		for rowStep := src.Row + 1; rowStep < dst.Row; rowStep++ {
			if b.board[rowStep][colStep] != nil {
				return false
			}
			colStep++
		}
		return true
	} else if (dst.Col-src.Col) == (dst.Row-src.Row) && (dst.Col-src.Col) < 0 {
		// Move down-left
		colStep := src.Col - 1
		for rowStep := src.Row - 1; rowStep > dst.Row; rowStep-- {
			if b.board[rowStep][colStep] != nil {
				return false
			}
			colStep--
		}
	} else if (dst.Col-src.Col) == (dst.Row-src.Row)*-1 && (dst.Col-src.Col) < 0 {
		// Move up-left
		rowStep := src.Row + 1
		for colStep := src.Col - 1; colStep > dst.Col; colStep-- {
			if b.board[rowStep][colStep] != nil {
				return false
			}
			rowStep++
		}
		return true
	} else if (dst.Col-src.Col) == (dst.Row-src.Row)*-1 && (dst.Col-src.Col) > 0 {
		// Move down-right
		rowStep := src.Row - 1
		for colStep := src.Col + 1; colStep < dst.Col; colStep++ {
			if b.board[src.Col][colStep] != nil {
				return false
			}
			rowStep--
		}
		return true
	}
	return false
}

func (b *standardBoard) isValidKnightMove(piece *types.Piece, dst types.Location) bool {
	src := piece.Location
	return dst.Row == src.Row+2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Row == src.Row-2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Col == src.Col+2 && math.Abs(float64(dst.Row-src.Row)) == 1 ||
		dst.Col == src.Col-2 && math.Abs(float64(dst.Row-src.Row)) == 1
}

func (b *standardBoard) isValidPawnMove(piece *types.Piece, dst types.Location) bool {
	src := piece.Location
	if piece.Color == types.White {
		if src.Col == dst.Col && dst.Row == src.Row+1 && b.board[dst.Row][dst.Col] == nil {
			return true
		} else if (dst.Col == src.Col+1 || dst.Col == src.Col-1) && dst.Row == src.Row+1 && b.board[dst.Row][dst.Col] != nil {
			return true
		}
	} else {
		if src.Col == dst.Col && dst.Row == src.Row-1 && b.board[dst.Row][dst.Col] == nil {
			return true
		} else if (dst.Col == src.Col-1 || dst.Col == src.Col+1) && dst.Row == src.Row-1 && b.board[dst.Row][dst.Col] != nil {
			return true
		}
	}
	return false
}

func (b *standardBoard) isValidMove(src, dst types.Location) bool {
	piece := b.board[src.Row][src.Col]
	if piece == nil {
		return false
	}
	switch piece.Type {
	case types.King:
		return b.isValidKingMove(piece, dst)
	case types.Rook:
		return b.isValidRookMove(piece, dst)
	case types.Pawn:
		return b.isValidPawnMove(piece, dst)
	case types.Bishop:
		return b.isValidBishopMove(piece, dst)
	case types.Queen:
		return b.isValidBishopMove(piece, dst) || b.isValidRookMove(piece, dst)
	case types.Knight:
		return b.isValidKnightMove(piece, dst)
	default:
		return false
	}
}

func newStandardBoard() *standardBoard {
	kings := map[types.Color]*types.Piece{
		types.White: {
			Type:     types.King,
			Color:    types.White,
			Location: types.Location{Row: 0, Col: 4},
		},
		types.Black: {
			Type:     types.King,
			Color:    types.Black,
			Location: types.Location{Row: 7, Col: 4},
		},
	}

	return &standardBoard{
		kings: kings,
		board: [8][8]*types.Piece{
			{
				{
					Type:     types.Rook,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 0},
				},
				{
					Type:     types.Knight,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 1},
				},
				{
					Type:     types.Bishop,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 2},
				},
				{
					Type:     types.Queen,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 3},
				},
				kings[types.White],
				{
					Type:     types.Bishop,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 5},
				},
				{
					Type:     types.Knight,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 6},
				},
				{
					Type:     types.Rook,
					Color:    types.White,
					Location: types.Location{Row: 0, Col: 7},
				},
			},
			{
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 0},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 1},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 2},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 3},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 4},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 5},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 6},
				},
				{
					Type:     types.Pawn,
					Color:    types.White,
					Location: types.Location{Row: 1, Col: 7},
				},
			},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 0},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 1},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 2},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 3},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 4},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 5},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 6},
				},
				{
					Type:     types.Pawn,
					Color:    types.Black,
					Location: types.Location{Row: 6, Col: 7},
				},
			},
			{
				{
					Type:     types.Rook,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 0},
				},
				{
					Type:     types.Knight,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 1},
				},
				{
					Type:     types.Bishop,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 2},
				},
				{
					Type:     types.Queen,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 3},
				},
				kings[types.Black],
				{
					Type:     types.Bishop,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 5},
				},
				{
					Type:     types.Knight,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 6},
				},
				{
					Type:     types.Rook,
					Color:    types.Black,
					Location: types.Location{Row: 7, Col: 7},
				},
			},
		},
	}
}
