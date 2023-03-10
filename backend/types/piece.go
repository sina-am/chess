package types

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

const (
	White Color = iota
	Black
)

type Piece struct {
	Type     PieceType
	Color    Color
	Location Location
	IsDead   bool
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
