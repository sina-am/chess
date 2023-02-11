package engine

import (
	"context"
	"fmt"

	"github.com/sina-am/chess/types"
)

type standardGame struct {
	board   Board
	players [2]*types.Player
	kings   [2]*types.Piece
}

func NewStandardGame(board Board, players [2]*types.Player) (*standardGame, error) {
	return &standardGame{
		board:   board,
		players: players,
	}, nil
}

func (g *standardGame) Play(ctx context.Context, player *types.Player, src, dst types.Location) error {
	if !player.Turn {
		return fmt.Errorf("not your turn")
	}

	piece, err := g.board.GetPiece(src)
	if err != nil {
		return err
	}

	if piece.Color != player.Color {
		return fmt.Errorf("this is not the player's piece")
	}

	return g.board.MovePiece(piece, dst)
}

// func (g *standardGame) getOpponent() *types.Player {
// 	if g.Turn.Color == types.White {
// 		return g.Players[1]
// 	} else {
// 		return g.Players[0]
// 	}
// }
// func (g *standardGame) findKingLocation(color types.Color) Location {
// 	// TODO: use efficient algorithm
// 	for i := 0; i < 8; i++ {
// 		for j := 0; j < 8; j++ {
// 			if g.Board[i][j] != nil {
// 				if g.Board[i][j].Color == color && g.Board[i][j].Type == types.King {
// 					return Location{row: i, col: j}
// 				}
// 			}
// 		}
// 	}
// 	panic("can't find the king")
// }

// func (g *standardGame) changeTurn() {
// 	g.Turn = g.getOpponent()
// }

// func (g *standardGame) checkCheck(piece *types.Piece, loc Location) bool {
// 	opponent := g.getOpponent()
// 	kingLocation := g.findKingLocation(opponent.Color)
// 	if g.isValidMove(piece, loc, kingLocation) {
// 		return true
// 	} else {
// 		return false
// 	}
// }
