package chess

import (
	"fmt"
	"math"
)

type Reason string

const (
	Checkmate Reason = "checkmate"
	Stalemate Reason = "stalemate"
	Timeout   Reason = "timeout"
	Abandoned Reason = "abandoned"
	Resign    Reason = "resign"
	Draw      Reason = "draw"
)

type Result struct {
	Reason      Reason
	WinnerColor Color
}

var NoResult = Result{}

type castling struct {
	left  bool
	right bool
}

type ChessEngine struct {
	board        [8][8]*Piece
	kings        map[Color]*Piece
	pieces       map[Color][]*Piece
	castleRights map[Color]castling
	turn         Color

	possibleMoves map[*Piece][]Location
	finished      bool
	result        Result
}

func NewEngine() *ChessEngine {
	engine := &ChessEngine{
		board:         newStandardBoard(),
		kings:         make(map[Color]*Piece, 2),
		pieces:        map[Color][]*Piece{White: {}, Black: {}},
		turn:          White,
		possibleMoves: map[*Piece][]Location{},
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
		result: NoResult,
	}
	engine.kings = map[Color]*Piece{
		White: engine.board[0][4],
		Black: engine.board[7][4],
	}
	engine.pieces[White] = append(engine.pieces[White], engine.board[0][:]...)
	engine.pieces[White] = append(engine.pieces[White], engine.board[1][:]...)
	engine.pieces[Black] = append(engine.pieces[Black], engine.board[7][:]...)
	engine.pieces[Black] = append(engine.pieces[Black], engine.board[6][:]...)

	engine.generatePossibleMoves()
	return engine
}

func (g *ChessEngine) GetBoard() [8][8]*Piece {
	return g.board
}

func NewFromPieces(pieces []*Piece) *ChessEngine {
	engine := &ChessEngine{
		board:         newBoardFromPieces(pieces),
		kings:         map[Color]*Piece{White: nil, Black: nil},
		pieces:        map[Color][]*Piece{White: {}, Black: {}},
		turn:          White,
		possibleMoves: map[*Piece][]Location{},
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

	engine.generatePossibleMoves()
	return engine
}

func (g *ChessEngine) GetResult() Result {
	if g.finished {
		return g.result
	}
	return NoResult
}
func (g *ChessEngine) generatePossibleMoves() {
	pieces := g.pieces[g.turn]

	for i := range pieces {
		if pieces[i].Captured {
			continue
		}
		switch pieces[i].Type {
		case King:
			g.possibleMoves[pieces[i]] = g.generateKingPossibleMoves(pieces[i])
		case Pawn:
			g.possibleMoves[pieces[i]] = g.generatePawnPossibleMoves(pieces[i])
		case Bishop:
			g.possibleMoves[pieces[i]] = g.generateBishopPossibleMoves(pieces[i])
		case Rook:
			g.possibleMoves[pieces[i]] = g.generateRookPossibleMoves(pieces[i])
		case Knight:
			g.possibleMoves[pieces[i]] = g.generateKnightPossibleMoves(pieces[i])
		case Queen:
			g.possibleMoves[pieces[i]] = append(g.generateBishopPossibleMoves(pieces[i]), g.generateRookPossibleMoves(pieces[i])...)
		}
	}
}
func (g *ChessEngine) occupiable(piece *Piece, loc Location) bool {
	if err := loc.Validate(); err != nil {
		return false
	}
	if (g.board[loc.Row][loc.Col] != nil) && (g.board[loc.Row][loc.Col].Color == piece.Color) {
		return false
	}

	rb := NewRollBack(g)
	rb.Do(Move{From: piece.Location, To: loc})
	isChecked := g.isChecked(piece.Color)
	rb.RollBack()

	return !isChecked
}

func (g *ChessEngine) generateRookPossibleMoves(piece *Piece) []Location {
	possibleMoves := []Location{}
	i := 0
	j := 0

	for i = piece.Location.Row + 1; i < 8; i++ {
		if g.occupiable(piece, Location{Row: i, Col: piece.Location.Col}) {
			possibleMoves = append(possibleMoves, Location{Row: i, Col: piece.Location.Col})
		}
	}
	for i = piece.Location.Row - 1; i >= 0; i-- {
		if g.occupiable(piece, Location{Row: i, Col: piece.Location.Col}) {
			possibleMoves = append(possibleMoves, Location{Row: i, Col: piece.Location.Col})
		}
	}
	for j = piece.Location.Col + 1; j < 8; j++ {
		if g.occupiable(piece, Location{Row: piece.Location.Row, Col: j}) {
			possibleMoves = append(possibleMoves, Location{Row: piece.Location.Row, Col: j})
		}
	}
	for j = piece.Location.Col - 1; j >= 0; j-- {
		if g.occupiable(piece, Location{Row: piece.Location.Row, Col: j}) {
			possibleMoves = append(possibleMoves, Location{Row: piece.Location.Row, Col: j})
		}
	}
	return possibleMoves
}

func (g *ChessEngine) generatePawnPossibleMoves(piece *Piece) []Location {
	possibleMoves := []Location{}

	if piece.Color == White {
		if g.board[piece.Location.Row+1][piece.Location.Col] == nil {
			possibleMoves = append(possibleMoves, Location{Row: piece.Location.Row + 1, Col: piece.Location.Col})
			// First move
			if piece.Location.Row == 1 && g.board[piece.Location.Row+2][piece.Location.Col] == nil {
				possibleMoves = append(possibleMoves, Location{Row: piece.Location.Row + 2, Col: piece.Location.Col})
			}
		}

		destinations := []Location{
			{Row: piece.Location.Row + 1, Col: piece.Location.Col + 1},
			{Row: piece.Location.Row + 1, Col: piece.Location.Col - 1},
		}
		for _, destination := range destinations {
			if err := destination.Validate(); err != nil {
				continue
			}
			if (g.board[destination.Row][destination.Col] != nil) && (g.board[destination.Row][destination.Col].Color != piece.Color) {
				possibleMoves = append(possibleMoves, destination)
			}
		}
	} else {
		if g.board[piece.Location.Row-1][piece.Location.Col] == nil {
			possibleMoves = append(possibleMoves, Location{Row: piece.Location.Row - 1, Col: piece.Location.Col})
			// First move
			if piece.Location.Row == 6 && g.board[piece.Location.Row-2][piece.Location.Col] == nil {
				possibleMoves = append(possibleMoves, Location{Row: piece.Location.Row - 2, Col: piece.Location.Col})
			}
		}
		destinations := []Location{
			{Row: piece.Location.Row - 1, Col: piece.Location.Col + 1},
			{Row: piece.Location.Row - 1, Col: piece.Location.Col - 1},
		}
		for _, destination := range destinations {
			if err := destination.Validate(); err != nil {
				continue
			}
			if (g.board[destination.Row][destination.Col] != nil) && (g.board[destination.Row][destination.Col].Color != piece.Color) {
				possibleMoves = append(possibleMoves, destination)
			}
		}
	}

	finalPossibleMoves := []Location{}
	for _, location := range possibleMoves {
		if g.occupiable(piece, location) {
			finalPossibleMoves = append(finalPossibleMoves, location)
		}
	}
	return finalPossibleMoves
}

func (g *ChessEngine) generateBishopPossibleMoves(piece *Piece) []Location {
	possibleMoves := []Location{}

	j := piece.Location.Col + 1
	// Down-right
	for i := piece.Location.Row + 1; i < 8; i++ {
		if g.occupiable(piece, Location{Row: i, Col: j}) {
			possibleMoves = append(possibleMoves, Location{Row: i, Col: j})
			j += 1
		}
	}
	j = piece.Location.Col - 1
	for i := piece.Location.Row + 1; i < 8; i++ {
		// Down-left
		if g.occupiable(piece, Location{Row: i, Col: j}) {
			possibleMoves = append(possibleMoves, Location{Row: i, Col: j})
			j -= 1
		}
	}
	j = piece.Location.Col + 1
	for i := piece.Location.Row - 1; i >= 0; i-- {
		// Up-right
		if g.occupiable(piece, Location{Row: i, Col: j}) {
			possibleMoves = append(possibleMoves, Location{Row: i, Col: j})
			j += 1
		}
	}
	j = piece.Location.Col - 1
	for i := piece.Location.Row - 1; i >= 0; i-- {
		// Up-left
		if g.occupiable(piece, Location{Row: i, Col: j}) {
			possibleMoves = append(possibleMoves, Location{Row: i, Col: j})
			j -= 1
		}
	}
	return possibleMoves
}

func (g *ChessEngine) generateKingPossibleMoves(king *Piece) []Location {
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
	possibleMoves := []Location{}
	for _, location := range locations {
		if !g.occupiable(king, location) {
			continue
		}
		possibleMoves = append(possibleMoves, location)
	}

	if g.castleRights[king.Color].left {
		if (g.board[king.Location.Row][3] == nil) && (g.board[king.Location.Row][2] == nil) && (g.board[king.Location.Row][1] == nil) {
			possibleMoves = append(possibleMoves, Location{Row: king.Location.Row, Col: king.Location.Col - 2})
		}
	}
	if g.castleRights[king.Color].right {
		if (g.board[king.Location.Row][5] == nil) && (g.board[king.Location.Row][6] == nil) {
			possibleMoves = append(possibleMoves, Location{Row: king.Location.Row, Col: king.Location.Col + 2})
		}
	}
	return possibleMoves
}

func (g *ChessEngine) generateKnightPossibleMoves(piece *Piece) []Location {
	possibleMoves := []Location{}
	destinations := []Location{
		{Row: piece.Location.Row + 2, Col: piece.Location.Col + 1},
		{Row: piece.Location.Row + 2, Col: piece.Location.Col - 1},
		{Row: piece.Location.Row - 2, Col: piece.Location.Col + 1},
		{Row: piece.Location.Row - 2, Col: piece.Location.Col - 1},
		{Row: piece.Location.Row + 1, Col: piece.Location.Col + 2},
		{Row: piece.Location.Row + 1, Col: piece.Location.Col - 2},
		{Row: piece.Location.Row - 1, Col: piece.Location.Col + 2},
		{Row: piece.Location.Row - 1, Col: piece.Location.Col - 2},
	}

	for _, location := range destinations {
		if !g.occupiable(piece, location) {
			continue
		}
		possibleMoves = append(possibleMoves, location)
	}

	return possibleMoves
}

func (g *ChessEngine) finish(reason Reason, winner Color) {
	g.result = Result{
		Reason:      reason,
		WinnerColor: winner,
	}
	g.finished = true
}

func (g *ChessEngine) IsInPossibleMoves(piece *Piece, loc Location) bool {
	possibleMoves := g.possibleMoves[piece]
	for i := range possibleMoves {
		if possibleMoves[i].Equals(loc) {
			return true
		}
	}
	return false
}

func (g *ChessEngine) Play(playerColor Color, move Move) error {
	if g.finished {
		return ErrGameEnd
	}

	if err := move.Validate(); err != nil {
		return err
	}

	if playerColor != g.turn {
		return ErrNotPlayersTurn
	}

	piece := g.board[move.From.Row][move.From.Col]
	if piece == nil {
		return ErrInvalidPieceMove
	}
	if piece.Color != playerColor {
		return ErrInvalidPieceMove
	}

	if !g.IsInPossibleMoves(piece, move.To) {
		return ErrInvalidPieceMove
	}

	rb := NewRollBack(g)
	rb.Do(move)

	g.switchTurn()

	if result := g.checkResult(); result != NoResult {
		g.finish(result.Reason, result.WinnerColor)
	}
	return nil
}

func (g *ChessEngine) checkResult() Result {
	for _, locations := range g.possibleMoves {
		if len(locations) != 0 {
			return NoResult
		}
	}
	king := g.kings[g.turn]
	if g.isChecked(king.Color) {
		return Result{
			Reason:      Checkmate,
			WinnerColor: king.Color.OppositeColor(),
		}
	}
	return Result{
		Reason:      Stalemate,
		WinnerColor: Empty,
	}
}

func (g *ChessEngine) Print() {
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

func (g *ChessEngine) switchTurn() {
	if g.turn == White {
		g.turn = Black
	} else {
		g.turn = White
	}
	g.generatePossibleMoves()
}
func (g *ChessEngine) SwitchTurn() {
	if g.turn == White {
		g.turn = Black
	} else {
		g.turn = White
	}
	g.generatePossibleMoves()
}

func (g *ChessEngine) isChecked(color Color) bool {
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

func (g *ChessEngine) isValidMove(src, dst Location) bool {
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

func (g *ChessEngine) isValidCastling(src, dst Location) bool {
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

func (g *ChessEngine) isValidKingMove(src, dst Location) bool {
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

func (g *ChessEngine) isValidRookMove(src, dst Location) bool {
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

func (g *ChessEngine) isValidBishopMove(src, dst Location) bool {
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

func (g *ChessEngine) isValidKnightMove(src, dst Location) bool {
	return dst.Row == src.Row+2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Row == src.Row-2 && math.Abs(float64(dst.Col-src.Col)) == 1 ||
		dst.Col == src.Col+2 && math.Abs(float64(dst.Row-src.Row)) == 1 ||
		dst.Col == src.Col-2 && math.Abs(float64(dst.Row-src.Row)) == 1
}

func (g *ChessEngine) isValidPawnMove(src, dst Location) bool {
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
