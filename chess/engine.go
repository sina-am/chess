package chess

import (
	"fmt"
	"math"
)

type castling struct {
	left  bool
	right bool
}

type chessEngine struct {
	board        [8][8]*Piece
	kings        map[Color]*Piece
	pieces       map[Color][]*Piece
	castleRights map[Color]castling
	turn         Color
	finished     bool
	winner       Color
}

func NewEngine() *chessEngine {
	engine := &chessEngine{
		board:  newStandardBoard(),
		kings:  make(map[Color]*Piece, 2),
		pieces: map[Color][]*Piece{White: {}, Black: {}},
		turn:   White,
		castleRights: map[Color]castling{
			White: {
				left:  true,
				right: true,
			},
			Black: {
				left:  true,
				right: true,
			},
		},
	}
	engine.kings = map[Color]*Piece{
		White: engine.board[0][4],
		Black: engine.board[7][4],
	}
	engine.pieces[White] = append(engine.pieces[White], engine.board[0][:]...)
	engine.pieces[White] = append(engine.pieces[White], engine.board[1][:]...)
	engine.pieces[Black] = append(engine.pieces[Black], engine.board[7][:]...)
	engine.pieces[Black] = append(engine.pieces[Black], engine.board[6][:]...)

	return engine
}
func NewFromPieces(pieces []*Piece) *chessEngine {
	engine := &chessEngine{
		board:  newBoardFromPieces(pieces),
		kings:  map[Color]*Piece{White: nil, Black: nil},
		pieces: map[Color][]*Piece{White: {}, Black: {}},
		turn:   White,
		castleRights: map[Color]castling{
			White: {
				left:  true,
				right: true,
			},
			Black: {
				left:  true,
				right: true,
			},
		},
	}

	for _, piece := range pieces {
		if piece.Type == King {
			engine.kings[piece.Color] = piece
		}
		engine.pieces[piece.Color] = append(engine.pieces[piece.Color], piece)
	}
	return engine
}

func (g *chessEngine) IsFinished() bool {
	return g.finished
}

func (g *chessEngine) Play(playerColor Color, move Move) error {
	if g.finished {
		return ErrGameEnd
	}

	if err := move.Validate(); err != nil {
		return err
	}

	if playerColor != g.turn {
		return ErrNotPlayersTurn
	}

	if !g.isValidMove(move.From, move.To) {
		return ErrInvalidPieceMove
	}
	if g.board[move.From.Row][move.From.Col] == nil {
		return ErrInvalidPieceMove
	}
	if g.board[move.From.Row][move.From.Col].Color != playerColor {
		return ErrInvalidPieceMove
	}
	if err := g.movePiece(playerColor, move); err != nil {
		return err
	}
	g.switchTurn()

	if g.IsCheckmated() {
		g.winner = g.turn.OppositeColor()
		g.finished = true
	}
	return nil
}

type RollBackMovement struct {
	game          *chessEngine
	capturedPiece *Piece
	move          Move
	castled       bool
	rolledBack    bool
	promoted      bool
}

func NewRollBack(game *chessEngine) *RollBackMovement {
	return &RollBackMovement{
		game: game,
	}
}

func (r *RollBackMovement) CheckPromotion(move Move) {
	piece := r.game.board[move.From.Row][move.From.Col]
	if piece.Type == Pawn {
		if piece.Color == Black && move.To.Row == 0 {
			piece.Type = Queen
			r.promoted = true
		} else if move.To.Row == 7 {
			piece.Type = Queen
			r.promoted = true
		}
	}
}
func (r *RollBackMovement) isCastling() bool {
	if (r.game.board[r.move.From.Row][r.move.From.Col].Type == King) && (math.Abs(float64(r.move.From.Col)-float64(r.move.To.Col)) == 2) {
		return true
	}
	return false
}

func (r *RollBackMovement) doCastling() {
	color := r.game.board[r.move.From.Row][r.move.From.Col].Color
	backRow := 0
	if color == Black {
		backRow = 7
	}

	if r.move.To.Col == 6 {
		r.game.castleRights[color] = castling{
			left:  r.game.castleRights[color].left,
			right: false,
		}
		r.game.board[backRow][5] = r.game.board[backRow][7]
		r.game.board[backRow][6] = r.game.board[backRow][4]
		r.game.board[backRow][4] = nil
		r.game.board[backRow][7] = nil
		return
	}
	if r.move.To.Col == 2 {
		r.game.castleRights[color] = castling{
			left:  false,
			right: r.game.castleRights[color].right,
		}
		r.game.board[backRow][3] = r.game.board[backRow][0]
		r.game.board[backRow][2] = r.game.board[backRow][4]
		r.game.board[backRow][4] = nil
		r.game.board[backRow][0] = nil
	}
}

func (r *RollBackMovement) Do(move Move) {
	r.move = move

	if r.isCastling() {
		r.doCastling()
		r.castled = true
		return
	}

	if r.game.board[move.To.Row][move.To.Col] != nil {
		r.capturedPiece = r.game.board[move.To.Row][move.To.Col]
		r.capturedPiece.Captured = true
	}

	r.CheckPromotion(move)

	r.game.board[move.To.Row][move.To.Col] = r.game.board[move.From.Row][move.From.Col]
	r.game.board[move.From.Row][move.From.Col] = nil
	r.game.board[move.To.Row][move.To.Col].Location.Col = move.To.Col
	r.game.board[move.To.Row][move.To.Col].Location.Row = move.To.Row
}

func (r *RollBackMovement) rollBackCastle() {
	color := r.game.board[r.move.To.Row][r.move.To.Col].Color

	r.game.board[r.move.To.Row][4] = r.game.board[r.move.To.Row][r.move.To.Col]
	r.game.board[r.move.To.Row][r.move.To.Col] = nil
	if r.move.To.Col == 2 {
		r.game.board[r.move.To.Row][0] = r.game.board[r.move.From.Row][3]
		r.game.board[r.move.To.Row][3] = nil
		r.game.castleRights[color] = castling{
			left:  true,
			right: r.game.castleRights[color].right,
		}
	} else {
		r.game.board[r.move.To.Row][7] = r.game.board[r.move.To.Row][5]
		r.game.board[r.move.To.Row][5] = nil
		r.game.castleRights[color] = castling{
			left:  r.game.castleRights[color].left,
			right: true,
		}
	}
}

func (r *RollBackMovement) RollBack() {
	if r.rolledBack {
		return
	}

	if r.castled {
		r.rollBackCastle()
		r.rolledBack = true
		return
	}
	if r.capturedPiece != nil {
		r.capturedPiece.Captured = false
	}

	r.game.board[r.move.From.Row][r.move.From.Col] = r.game.board[r.move.To.Row][r.move.To.Col]
	r.game.board[r.move.To.Row][r.move.To.Col] = r.capturedPiece
	r.game.board[r.move.From.Row][r.move.From.Col].Location.Col = r.move.From.Col
	r.game.board[r.move.From.Row][r.move.From.Col].Location.Row = r.move.From.Row

	if r.promoted {
		r.game.board[r.move.From.Row][r.move.From.Col].Type = Pawn
	}
	r.rolledBack = true
}

func (g *chessEngine) IsCheckmated() bool {
	king := g.kings[g.turn]
	if !g.isChecked(king.Color) {
		return false
	}

	locations := []Location{
		{king.Location.Row, king.Location.Col - 1},
		{king.Location.Row, king.Location.Col + 1},
		{king.Location.Row + 1, king.Location.Col},
		{king.Location.Row - 1, king.Location.Col},
		{king.Location.Row + 1, king.Location.Col + 1},
		{king.Location.Row + 1, king.Location.Col - 1},
		{king.Location.Row - 1, king.Location.Col + 1},
		{king.Location.Row - 1, king.Location.Col - 1},
	}
	for _, location := range locations {
		if err := location.Validate(); err != nil {
			continue
		}
		if g.isValidMove(king.Location, location) {
			rb := NewRollBack(g)

			rb.Do(Move{From: king.Location, To: location})

			isChecked := g.isChecked(king.Color)

			rb.RollBack()

			if !isChecked {
				return false
			}
		}
	}
	return true
}

func (g *chessEngine) GetWinner() Color {
	if g.finished {
		return g.winner
	}
	return Empty
}

func (g *chessEngine) Print() {
	fmt.Println("########")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if g.board[i][j] != nil {
				fmt.Printf("%s ", g.board[i][j].String())
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
	fmt.Println("########")
}

func (g *chessEngine) movePiece(playerColor Color, move Move) error {
	rb := NewRollBack(g)

	rb.Do(move)
	if g.isChecked(playerColor) {
		rb.RollBack()
		return ErrChecked
	}

	return nil
}
func (g *chessEngine) switchTurn() {
	if g.turn == White {
		g.turn = Black
	} else {
		g.turn = White
	}
}

func (g *chessEngine) isChecked(color Color) bool {
	king := g.kings[color]
	if king == nil {
		panic("king is not there")
	}

	opponentPieces := g.pieces[king.Color.OppositeColor()]
	for _, piece := range opponentPieces {
		if (!piece.Captured) && g.isValidMove(piece.Location, king.Location) {
			return true
		}
	}
	return false
}

func (g *chessEngine) isValidMove(src, dst Location) bool {
	piece := g.board[src.Row][src.Col]

	if piece == nil {
		return false
	}

	if g.board[dst.Row][dst.Col] != nil {
		if g.board[src.Row][src.Col].Color == g.board[dst.Row][dst.Col].Color {
			return false
		}
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

func (g *chessEngine) isValidCastling(src, dst Location) bool {
	color := g.board[src.Row][src.Col].Color
	backRow := 0
	if color == Black {
		backRow = 7
	}
	if src.Row != backRow || dst.Row != backRow || src.Col != 4 {
		return false
	}
	if src.Col == 4 && dst.Col == 6 {
		// Left castle
		if g.board[backRow][5] == nil && g.board[backRow][6] == nil {
			if g.castleRights[color].left {
				return true
			} else {
				return false
			}
		} else {
			return false
		}
	} else if src.Col == 4 && dst.Col == 2 {
		// Right castle
		if g.board[backRow][3] == nil && g.board[backRow][2] == nil && g.board[backRow][1] == nil {
			if g.castleRights[color].right {
				return true
			} else {
				return false
			}
		}
		return false
	}
	return false
}

func (g *chessEngine) isValidKingMove(src, dst Location) bool {
	// Castling
	if math.Abs(float64(src.Col)-float64(dst.Col)) == 2 && math.Abs(float64(src.Row)-float64(dst.Row)) == 0 {
		return g.isValidCastling(src, dst)
	}

	return (src.Col == dst.Col &&
		math.Abs(float64(src.Row)-float64(dst.Row)) == 1) ||
		(src.Row == dst.Row &&
			math.Abs(float64(src.Col)-float64(dst.Col)) == 1) ||
		(src.Col == dst.Col+1) && (src.Row == dst.Row+1) ||
		(src.Col == dst.Col+1) && (src.Row == dst.Row-1) ||
		(src.Col == dst.Col-1) && (src.Row == dst.Row+1) ||
		(src.Col == dst.Col-1) && (src.Row == dst.Row-1)
}

func (g *chessEngine) isValidRookMove(src, dst Location) bool {
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
		for colStep := src.Col - 1; colStep > dst.Col; colStep-- {
			if g.board[src.Row][colStep] != nil {
				return false
			}
		}
		return true
	}
	return false
}

func (g *chessEngine) isValidBishopMove(src, dst Location) bool {
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

func (g *chessEngine) isValidKnightMove(src, dst Location) bool {
	return dst.Row == src.Row+2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Row == src.Row-2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Col == src.Col+2 && math.Abs(float64(dst.Row-src.Row)) == 1 ||
		dst.Col == src.Col-2 && math.Abs(float64(dst.Row-src.Row)) == 1
}

func (g *chessEngine) isValidPawnMove(src, dst Location) bool {
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

func newBoardFromPieces(pieces []*Piece) [8][8]*Piece {
	var board [8][8]*Piece
	for _, piece := range pieces {
		board[piece.Location.Row][piece.Location.Col] = piece
	}
	return board
}

func newStandardBoard() [8][8]*Piece {
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
