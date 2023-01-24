package types

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

type PieceNumber int

type Piece struct {
	Type        PieceType
	PieceNumber PieceNumber
	Color       Color
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
