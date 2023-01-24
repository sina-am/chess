package engine

import (
	"context"
	"time"

	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/types"
)

type Game interface {
	Request(ctx context.Context, from *types.User, to *types.User) error
	Accept(ctx context.Context, user *types.User, game *types.Game) error
	// MovePiece(user *types.User, game *types.Game, src types.Location, dst types.Location) error
}

type standardGame struct {
	Database database.Database
}

func NewStandardGame(db database.Database) (*standardGame, error) {
	return &standardGame{
		Database: db,
	}, nil
}

func (g *standardGame) Request(ctx context.Context, from *types.User, to *types.User) error {
	newGame := &types.Game{
		StartedAt: time.Time{},
		Duration:  10,
		Players: []types.Player{
			{
				UserId:       from.Id,
				Color:        types.White,
				IsChecked:    false,
				IsCheckmated: false,
				Turn:         true,
			},
			{
				UserId:       to.Id,
				Color:        types.Black,
				IsChecked:    false,
				IsCheckmated: false,
				Turn:         false,
			},
		},
		StartedBy:  from.Id,
		IsAccepted: false,
	}
	return g.Database.InsertGame(ctx, newGame)
}

func (g *standardGame) Accept(ctx context.Context, user *types.User, game *types.Game) error {
	game.Pieces = makePieces()
	game.Board = makeBoard(game.Pieces)
	game.IsAccepted = true
	return g.Database.UpdateUserGame(ctx, game)
}

// 	pieces := makePieces()
// 	board := makeBoard(pieces)

// 	return &standardGame{
// 		Players:    players,
// 		Pieces:     pieces,
// 		Board:      board,
// 		TookPieces: make([]types.Piece, 0),
// 		Turn:       players[0],
// 	}
// }

// func (g *standardGame) MovePiece(playerColor types.Color, src Location, dst Location) error {
// 	if playerColor != g.Turn.Color {
// 		return fmt.Errorf("not your turn")
// 	}

// 	piece := g.Board[src.row][src.col]
// 	if piece == nil {
// 		return fmt.Errorf("invalid source location")
// 	}
// 	if piece.Color != playerColor {
// 		return fmt.Errorf("this is not the player's piece")
// 	}
// 	if !g.isValidMove(piece, src, dst) {
// 		return fmt.Errorf("piece %s can't move from %s to %s", piece.String(), src.String(), dst.String())
// 	}

// 	return g.movePiece(piece, src, dst)
// }

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

// func (g *standardGame) movePiece(piece *types.Piece, src, dst Location) error {
// 	if g.Board[dst.row][dst.col] != nil {
// 		g.TookPieces = append(g.TookPieces, *g.Board[dst.row][dst.col])
// 	}
// 	g.Board[src.row][src.col] = nil
// 	g.Board[dst.row][dst.col] = piece

// 	if g.checkCheck(piece, dst) {
// 		opponent := g.getOpponent()
// 		opponent.IsChecked = true
// 	}
// 	g.changeTurn()
// 	return nil
// }
