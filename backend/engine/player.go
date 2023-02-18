package engine

import (
	"context"
	"fmt"

	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OnlinePlayer interface {
	WaitForStart(ctx context.Context) error
	WaitForMyTurn(ctx context.Context) error
	Play(ctx context.Context, src, dst types.Location) error
	GetGame() types.Game
}

type onlinePlayer struct {
	id    primitive.ObjectID
	game  Game
	evtCh chan gameEvent
}

func (p *onlinePlayer) WaitForStart(ctx context.Context) error {
	evt := <-p.evtCh
	if evt != startedEvent {
		return fmt.Errorf("something went wrong")
	}
	return nil
}

func (p *onlinePlayer) WaitForMyTurn(ctx context.Context) error {
	evt := <-p.evtCh
	if evt == yourTurnEvent {
		return nil
	}
	return fmt.Errorf("somthing went wrong")
}

func (p *onlinePlayer) Play(ctx context.Context, src, dst types.Location) error {
	return p.game.Play(ctx, p.id, src, dst)
}

func (p *onlinePlayer) GetGame() types.Game {
	return p.game.GetGame()
}
