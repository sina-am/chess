package chess

import "fmt"

type Location struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (loc Location) Validate() error {
	if loc.Row >= 0 && loc.Row < 8 && loc.Col >= 0 && loc.Col < 8 {
		return nil
	}
	return ErrOutOfBoardMove
}

func (loc *Location) String() string {
	return fmt.Sprintf("(%d, %d)", loc.Row, loc.Col)
}

func (loc *Location) Equals(loc2 Location) bool {
	return loc.Row == loc2.Row && loc.Col == loc2.Col
}

type PieceType int

const (
	King PieceType = iota
	Rook
	Bishop
	Queen
	Knight
	Pawn
)

type Color int

func (c Color) String() string {
	if c == White {
		return "white"
	} else {
		return "black"
	}
}

const (
	White Color = iota
	Black
	Empty
)

func (color Color) OppositeColor() Color {
	if color == White {
		return Black
	} else if color == Black {
		return White
	} else {
		panic("invalid color")
	}
}

type Piece struct {
	Type     PieceType
	Color    Color
	Location Location
	Captured bool
}

func (b *Piece) String() string {
	switch b.Color {
	case Black:
		switch b.Type {
		case King:
			return "♔"
		case Rook:
			return "♖"
		case Bishop:
			return "♗"
		case Queen:
			return "♕"
		case Knight:
			return "♘"
		case Pawn:
			return "♙"
		default:
			panic("unknown piece")
		}
	case White:
		switch b.Type {
		case King:
			return "♚"
		case Rook:
			return "♜"
		case Bishop:
			return "♝"
		case Queen:
			return "♛"
		case Knight:
			return "♞"
		case Pawn:
			return "♟"
		default:
			panic("unknown piece")
		}
	default:
		panic("unknown color")
	}
}
