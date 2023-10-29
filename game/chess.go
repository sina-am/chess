package game

import (
	"errors"
	"fmt"
	"math"
)

var ErrChecked = errors.New("can't move there")
var ErrInvalidMove = errors.New("invalid move")

type gameEngine struct {
	board          [8][8]*Piece
	kings          map[Color]*Piece
	players        map[string]Color
	capturedPieces []*Piece
	turn           Color
}

func NewChess(players map[string]Color) *gameEngine {
	engine := &gameEngine{
		board:          NewStandardBoard(),
		kings:          make(map[Color]*Piece, 2),
		players:        players,
		turn:           White,
		capturedPieces: []*Piece{},
	}
	engine.kings = map[Color]*Piece{
		White: engine.board[0][4],
		Black: engine.board[7][4],
	}
	return engine
}
func NewChessFromPieces(players map[string]Color, pieces []*Piece, kings map[Color]*Piece) *gameEngine {
	engine := &gameEngine{
		board:          NewBoardFromPieces(pieces),
		kings:          kings,
		players:        players,
		turn:           White,
		capturedPieces: []*Piece{},
	}

	return engine
}
func (g *gameEngine) Print() {
	fmt.Println("########")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if g.board[i][j] != nil {
				fmt.Printf("%s ", g.board[i][j].String())
			} else {
				// First character is unicode U+0020
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
	fmt.Println("########")
}

func (g *gameEngine) switchTurn() {
	if g.turn == White {
		g.turn = Black
	} else {
		g.turn = White
	}
}

func (g *gameEngine) Exit(playerId string) (Color, error) {
	playerColor, found := g.players[playerId]
	if !found {
		return 0, fmt.Errorf("player with id %s not found", playerId)
	}

	return playerColor.OppositeColor(), nil
}

func (g *gameEngine) GetPlayers() []string {
	var ids []string
	for id := range g.players {
		ids = append(ids, id)
	}
	return ids
}

func (g *gameEngine) GetWinner() (string, error) {
	return "", fmt.Errorf("not implemented yet")
}

func (g *gameEngine) InGame(playerId string) bool {
	if _, found := g.players[playerId]; found {
		return true
	}
	return false
}

func (g *gameEngine) Play(playerId string, move Move) error {
	src := move.From
	dst := move.To

	playerColor, found := g.players[playerId]
	if !found {
		return fmt.Errorf("player not found")
	}
	if playerColor != g.turn {
		return fmt.Errorf("it's not %s turn", playerColor.String())
	}
	if playerColor != g.board[src.Row][src.Col].Color {
		return fmt.Errorf("it's not %s piece", playerColor.String())
	}

	if err := g.movePiece(playerColor, src, dst); err != nil {
		return err
	}
	// opponent := g.getOpponentPlayer()
	// if g.isChecked(opponent.Color) {
	// 	opponent.IsChecked = true
	// }
	g.switchTurn()
	return nil
}

func (g *gameEngine) isChecked(color Color) bool {
	king := g.kings[color]
	if king == nil {
		panic("king is not there")
	}
	opponentColor := king.Color.OppositeColor()
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if g.board[i][j] != nil && g.board[i][j].Color == opponentColor {
				if g.isValidMove(g.board[i][j].Location, king.Location) {
					return true
				}
			}
		}
	}
	return false
}
func (g *gameEngine) isValidMove(src, dst Location) bool {
	piece := g.board[src.Row][src.Col]
	if piece == nil {
		return false
	}
	switch piece.Type {
	case King:
		return g.isValidKingMove(src, dst)
	case Rook:
		return g.isValidRookMove(src, dst)
	case Pawn:
		return g.isValidPawnMove(src, dst)
	case Bishop:
		return g.isValidBishopMove(src, dst)
	case Queen:
		return g.isValidBishopMove(src, dst) || g.isValidRookMove(src, dst)
	case Knight:
		return g.isValidKnightMove(src, dst)
	default:
		return false
	}
}
func (g *gameEngine) movePiece(playerColor Color, src, dst Location) error {
	if !g.isValidMove(src, dst) {
		return ErrInvalidMove
	}

	var capturedPiece *Piece
	if g.board[dst.Row][dst.Col] != nil {
		capturedPiece = g.board[dst.Row][dst.Col]
	}

	g.board[dst.Row][dst.Col] = g.board[src.Row][src.Col]
	g.board[src.Row][src.Col] = nil

	g.board[dst.Row][dst.Col].Location.Col = dst.Col
	g.board[dst.Row][dst.Col].Location.Row = dst.Row

	// RoleBack
	if g.isChecked(playerColor) {
		g.board[src.Row][src.Col] = g.board[dst.Row][dst.Col]
		if capturedPiece != nil {
			g.board[dst.Row][dst.Col] = capturedPiece
		} else {
			g.board[dst.Row][dst.Col] = nil
		}

		g.board[src.Row][src.Col].Location.Col = src.Col
		g.board[src.Row][src.Col].Location.Row = src.Row
		return ErrChecked
	}

	g.capturedPieces = append(g.capturedPieces, g.board[dst.Row][dst.Col])
	return nil
}

func (g *gameEngine) isValidKingMove(src, dst Location) bool {
	return (src.Col == dst.Col &&
		math.Abs(float64(src.Row)-float64(dst.Row)) == 1) ||
		(src.Row == dst.Row &&
			math.Abs(float64(src.Col)-float64(dst.Col)) == 1) ||
		(src.Col == dst.Col+1) && (src.Row == dst.Row+1) ||
		(src.Col == dst.Col+1) && (src.Row == dst.Row-1) ||
		(src.Col == dst.Col-1) && (src.Row == dst.Row+1) ||
		(src.Col == dst.Col-1) && (src.Row == dst.Row-1)
}

func (g *gameEngine) isValidRookMove(src, dst Location) bool {
	if src.Col == dst.Col && src.Row < dst.Row {
		// Move down
		for rowStep := src.Row + 1; rowStep < dst.Row; rowStep++ {
			if g.board[rowStep][dst.Col] != nil {
				return false
			}
		}
		return true
	} else if src.Col == dst.Col && src.Row > dst.Row {
		// Move up
		for rowStep := src.Row - 1; rowStep > dst.Row; rowStep-- {
			if g.board[rowStep][dst.Col] != nil {
				return false
			}
		}
		return true
	} else if src.Row == dst.Row && src.Col < dst.Col {
		// Move right
		for colStep := src.Col + 1; colStep < dst.Col; colStep++ {
			if g.board[src.Row][colStep] != nil {
				return false
			}
		}
		return true
	} else if src.Row == dst.Row && src.Col > dst.Col {
		// Move left
		for colStep := src.Col - 1; colStep > dst.Col; colStep++ {
			if g.board[src.Row][colStep] != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (g *gameEngine) isValidBishopMove(src, dst Location) bool {
	if ((dst.Col - src.Col) == (dst.Row - src.Row)) && ((dst.Col - src.Col) > 0) {
		// Move up-right
		colStep := src.Col + 1
		for rowStep := src.Row + 1; rowStep < dst.Row; rowStep++ {
			if g.board[rowStep][colStep] != nil {
				return false
			}
			colStep++
		}
		return true
	} else if ((dst.Col - src.Col) == (dst.Row - src.Row)) && ((dst.Col - src.Col) < 0) {
		// Move down-left
		colStep := src.Col - 1
		for rowStep := src.Row - 1; rowStep > dst.Row; rowStep-- {
			if g.board[rowStep][colStep] != nil {
				return false
			}
			colStep--
		}
		return true
	} else if ((dst.Col - src.Col) == (dst.Row-src.Row)*(-1)) && ((dst.Col - src.Col) < 0) {
		// Move up-left
		rowStep := src.Row + 1
		for colStep := src.Col - 1; colStep > dst.Col; colStep-- {
			if g.board[rowStep][colStep] != nil {
				return false
			}
			rowStep++
		}
		return true
	} else if ((dst.Col - src.Col) == (dst.Row-src.Row)*(-1)) && ((dst.Col - src.Col) > 0) {
		// Move up-right
		rowStep := src.Row - 1
		for colStep := src.Col + 1; colStep < dst.Col; colStep++ {
			if g.board[rowStep][colStep] != nil {
				return false
			}
			rowStep--
		}
		return true
	}
	return false
}

func (g *gameEngine) isValidKnightMove(src, dst Location) bool {
	return dst.Row == src.Row+2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Row == src.Row-2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Col == src.Col+2 && math.Abs(float64(dst.Row-src.Row)) == 1 ||
		dst.Col == src.Col-2 && math.Abs(float64(dst.Row-src.Row)) == 1
}

func (g *gameEngine) isValidPawnMove(src, dst Location) bool {
	if g.board[src.Row][src.Col].Color == White {
		if src.Col == dst.Col && dst.Row == src.Row+1 && g.board[dst.Row][dst.Col] == nil {
			return true
		} else if src.Col == dst.Col && dst.Row == src.Row+2 && g.board[dst.Row][dst.Col] == nil && src.Row == 1 {
			return true
		} else if (dst.Col == src.Col+1 || dst.Col == src.Col-1) && dst.Row == src.Row+1 && g.board[dst.Row][dst.Col] != nil {
			return true
		}
	} else {
		if src.Col == dst.Col && dst.Row == src.Row-1 && g.board[dst.Row][dst.Col] == nil {
			return true
		} else if src.Col == dst.Col && dst.Row == src.Row-2 && g.board[dst.Row][dst.Col] == nil && src.Row == 6 {
			return true
		} else if (dst.Col == src.Col-1 || dst.Col == src.Col+1) && dst.Row == src.Row-1 && g.board[dst.Row][dst.Col] != nil {
			return true
		}
	}
	return false
}

func NewBoardFromPieces(pieces []*Piece) [8][8]*Piece {
	var board [8][8]*Piece
	for _, piece := range pieces {
		board[piece.Location.Row][piece.Location.Col] = piece
	}
	return board
}
func NewStandardBoard() [8][8]*Piece {
	return [8][8]*Piece{
		{
			{
				Type:     Rook,
				Color:    White,
				Location: Location{Row: 0, Col: 0},
			},
			{
				Type:     Knight,
				Color:    White,
				Location: Location{Row: 0, Col: 1},
			},
			{
				Type:     Bishop,
				Color:    White,
				Location: Location{Row: 0, Col: 2},
			},
			{
				Type:     Queen,
				Color:    White,
				Location: Location{Row: 0, Col: 3},
			},
			{
				Type:     King,
				Color:    White,
				Location: Location{Row: 0, Col: 4},
			},
			{
				Type:     Bishop,
				Color:    White,
				Location: Location{Row: 0, Col: 5},
			},
			{
				Type:     Knight,
				Color:    White,
				Location: Location{Row: 0, Col: 6},
			},
			{
				Type:     Rook,
				Color:    White,
				Location: Location{Row: 0, Col: 7},
			},
		},
		{
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 0},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 1},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 2},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 3},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 4},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 5},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 6},
			},
			{
				Type:     Pawn,
				Color:    White,
				Location: Location{Row: 1, Col: 7},
			},
		},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil, nil, nil, nil},
		{
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 0},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 1},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 2},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 3},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 4},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 5},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 6},
			},
			{
				Type:     Pawn,
				Color:    Black,
				Location: Location{Row: 6, Col: 7},
			},
		},
		{
			{
				Type:     Rook,
				Color:    Black,
				Location: Location{Row: 7, Col: 0},
			},
			{
				Type:     Knight,
				Color:    Black,
				Location: Location{Row: 7, Col: 1},
			},
			{
				Type:     Bishop,
				Color:    Black,
				Location: Location{Row: 7, Col: 2},
			},
			{
				Type:     Queen,
				Color:    Black,
				Location: Location{Row: 7, Col: 3},
			},
			{
				Type:     King,
				Color:    Black,
				Location: Location{Row: 7, Col: 4},
			},
			{
				Type:     Bishop,
				Color:    Black,
				Location: Location{Row: 7, Col: 5},
			},
			{
				Type:     Knight,
				Color:    Black,
				Location: Location{Row: 7, Col: 6},
			},
			{
				Type:     Rook,
				Color:    Black,
				Location: Location{Row: 7, Col: 7},
			},
		},
	}
}
