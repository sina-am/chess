package chess

import "fmt"

type Location struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Check if the location is out of the board
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

func (p PieceType) GetName() string {
	switch p {
	case King:
		return "king"
	case Rook:
		return "rook"
	case Bishop:
		return "bishop"
	case Queen:
		return "queen"
	case Knight:
		return "knight"
	case Pawn:
		return "pawn"
	default:
		panic("invalid piece type")
	}
}

type Color int

func (c Color) String() string {
	switch c {
	case White:
		return "white"
	case Black:
		return "black"
	case Empty:
		return "empty"
	default:
		panic("invalid color")
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
