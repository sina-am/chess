package game

import "fmt"

type Location struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func (loc *Location) String() string {
	return fmt.Sprintf("(%d, %d)", loc.Row, loc.Col)
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
)

func OppositeColor(color Color) Color {
	if color == White {
		return Black
	} else {
		return White
	}
}

type Piece struct {
	Type     PieceType
	Color    Color
	Location Location
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
			panic("unknown peice")
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
			panic("unknown peice")
		}
	default:
		panic("unkown color")
	}
}
