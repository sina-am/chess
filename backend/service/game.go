package service

import (
	"context"

	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/engine"
	"github.com/sina-am/chess/types"
)

type GameService interface {
	StartGame(ctx context.Context, user *types.User) (engine.OnlinePlayer, error)
}

type gameService struct {
	db    database.Database
	games []engine.Game
}

func NewGameService(db database.Database) (*gameService, error) {
	return &gameService{
		db:    db,
		games: []engine.Game{},
	}, nil
}

func (g *gameService) StartGame(ctx context.Context, user *types.User) (engine.OnlinePlayer, error) {
	for _, game := range g.games {
		if !game.HasStarted() {
			player, err := game.Join(ctx, user.Id)
			if err != nil {
				return nil, err
			}
			if err := player.WaitForStart(ctx); err != nil {
				return nil, err
			}
			return player, nil
		}
	}

	newGame, err := engine.NewOnlineGame(10)
	if err != nil {
		return nil, err
	}
	g.games = append(g.games, newGame)

	player, err := newGame.Join(ctx, user.Id)
	if err != nil {
		return nil, err
	}

	player.WaitForStart(ctx)
	return player, nil
}
