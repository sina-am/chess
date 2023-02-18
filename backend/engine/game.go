package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type gameEvent int

const (
	startedEvent gameEvent = iota
	yourTurnEvent
)

type Game interface {
	Play(ctx context.Context, userId primitive.ObjectID, src, dst types.Location) error
	Join(ctx context.Context, userId primitive.ObjectID) (*onlinePlayer, error)
	GetGame() types.Game
	HasStarted() bool
}

type onlineGame struct {
	game    types.Game
	board   Board
	players []*onlinePlayer
	started bool
	mu      sync.RWMutex
}

func NewOnlineGame(duration time.Duration) (*onlineGame, error) {
	return &onlineGame{
		game: types.Game{
			Id:       primitive.NewObjectID(),
			Duration: duration,
			Players:  make([]*types.Player, 0, 2),
		},
		board:   NewStandardBoard(),
		players: make([]*onlinePlayer, 0, 2),
		started: false,
	}, nil
}

func (g *onlineGame) GetGame() types.Game {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.game
}

func (g *onlineGame) Play(ctx context.Context, userId primitive.ObjectID, src, dst types.Location) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	player, err := g.getPlayer(userId)
	if err != nil {
		return err
	}

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

	if err := g.board.MovePiece(src, dst); err != nil {
		return err
	}

	g.switchTurn()
	return nil
}

func (g *onlineGame) addPlayer(userId primitive.ObjectID) (*onlinePlayer, error) {
	if len(g.game.Players) == 0 {
		g.game.Players = append(g.game.Players, &types.Player{
			UserId:       userId,
			IsChecked:    false,
			IsCheckmated: false,
			Color:        types.White,
			Turn:         true,
		})
	} else if len(g.game.Players) == 1 {
		g.game.Players = append(g.game.Players, &types.Player{
			UserId:       userId,
			IsChecked:    false,
			IsCheckmated: false,
			Color:        types.Black,
			Turn:         false,
		})
	} else {
		return nil, fmt.Errorf("two players already joined")
	}

	return &onlinePlayer{
		id:    userId,
		game:  g,
		evtCh: make(chan gameEvent, 1024),
	}, nil
}

func (g *onlineGame) Join(ctx context.Context, userId primitive.ObjectID) (*onlinePlayer, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.started {
		return nil, fmt.Errorf("game already started")
	}

	player, err := g.addPlayer(userId)
	if err != nil {
		return nil, err
	}
	g.players = append(g.players, player)

	if len(g.game.Players) == 2 {
		g.started = true
		g.announceStart()
	}
	return player, nil
}

func (g *onlineGame) HasStarted() bool {
	g.mu.RLock()
	defer g.mu.Unlock()

	return g.started
}

func (g *onlineGame) announceStart() {
	for i := range g.players {
		g.players[i].evtCh <- startedEvent
		if g.game.Players[i].Color == types.White {
			g.players[i].evtCh <- yourTurnEvent
		}
	}
	g.game.StartedAt = time.Now()
}

func (g *onlineGame) getPlayer(userId primitive.ObjectID) (*types.Player, error) {
	for i := range g.game.Players {
		if g.game.Players[i].UserId == userId {
			return g.game.Players[i], nil
		}
	}
	return nil, fmt.Errorf("invalid id")
}

func (g *onlineGame) switchTurn() {
	for i := range g.game.Players {
		if g.game.Players[i].Turn {
			g.game.Players[i].Turn = false
		} else {
			g.game.Players[i].Turn = true
			g.players[i].evtCh <- yourTurnEvent
		}
	}
}
