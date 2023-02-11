package service

import (
	"context"
	"time"

	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/types"
)

type GameService interface {
	Request(ctx context.Context, from *types.User, to *types.User) error
	Accept(ctx context.Context, user *types.User, game *types.Game) error
	Play(ctx context.Context, user *types.User, game *types.Game, src types.Location, dst types.Location) error
}

type gameService struct {
	db database.Database
}

func NewGameService(db database.Database) (*gameService, error) {
	return &gameService{
		db: db,
	}, nil
}

func (g *gameService) Request(ctx context.Context, duration time.Duration, from *types.User, to *types.User) error {
	newGame := &types.Game{
		Duration: duration,
		Players: []types.Player{
			{
				UserId: from.Id,
			},
			{
				UserId: to.Id,
			},
		},
		StartedBy:  from.Id,
		IsAccepted: false,
	}
	return g.db.InsertGame(ctx, newGame)
}

func (g *gameService) Accept(ctx context.Context, user *types.User, game *types.Game) error {
	
	game.IsAccepted = true
	return g.db.UpdateUserGame(ctx, game)
}
