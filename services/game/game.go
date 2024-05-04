package game

import (
	"context"
	"fmt"
	"time"

	"github.com/sina-am/chess/chess"
	"github.com/sina-am/chess/storage"
	"github.com/sina-am/chess/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClientEventType string

const (
	StartedClientEvent ClientEventType = "started"
	EndGameClientEvent ClientEventType = "ended"
	PlayedClientEvent  ClientEventType = "played"
)

type OnlineGame struct {
	Storage storage.Storage
	Players map[chess.Color]*onlinePlayer
	Game    chess.Chess

	drawOffered *onlinePlayer
}

type StartGamePayloadWSMsg struct {
	Opponent string `json:"opponent"`
	Tile     string `json:"tile"`
}
type StartGameWSMsg struct {
	Type    ClientEventType       `json:"type"`
	Payload StartGamePayloadWSMsg `json:"payload"`
}

func NewOnlineGame(s storage.Storage, p1, p2 *onlinePlayer, duration time.Duration) *OnlineGame {
	game := &OnlineGame{
		Storage: s,
		Players: map[chess.Color]*onlinePlayer{
			chess.White: p1,
			chess.Black: p2,
		},
		Game: chess.NewSession(duration),
	}

	p1.currentGame = game
	p1.status = StatusPlaying

	p2.currentGame = game
	p2.status = StatusPlaying

	p1.client.Send(StartGameWSMsg{
		Type: StartedClientEvent,
		Payload: StartGamePayloadWSMsg{
			Opponent: p2.user.GetName(),
			Tile:     chess.White.String(),
		},
	})
	p2.client.Send(StartGameWSMsg{
		Type: StartedClientEvent,
		Payload: StartGamePayloadWSMsg{
			Opponent: p1.user.GetName(),
			Tile:     chess.Black.String(),
		},
	})

	return game
}

func (g *OnlineGame) getPlayerColor(p *onlinePlayer) (chess.Color, error) {
	for color := range g.Players {
		if g.Players[color] == p {
			return color, nil
		}
	}
	return chess.Empty, fmt.Errorf("user is not in the game")
}

func (g *OnlineGame) Play(p *onlinePlayer, move chess.Move) error {
	color, err := g.getPlayerColor(p)
	if err != nil {
		return err
	}

	if err := g.Game.Play(color, move); err != nil {
		p.client.SendErr(err)
		return err
	}

	p.client.Send(map[string]string{
		"message": fmt.Sprintf("you played %d", move),
	})

	for _, pl := range g.Players {
		if p != pl {
			pl.client.Send(map[string]any{
				"type": PlayedClientEvent,
				"payload": map[string]any{
					"player": p.user.GetId(),
					"move":   move,
				},
			})
		}
	}

	if result := g.Game.GetResult(); result != chess.NoResult {
		return g.endGame(result)
	}
	return nil
}

func (g *OnlineGame) OfferDraw(p *onlinePlayer) {
	g.drawOffered = p
	opponent, _ := g.GetOpponentPlayer(p)
	opponent.client.Send(map[string]string{
		"type": "drawOffered",
	})
}

func (g *OnlineGame) RespondDraw(p *onlinePlayer, accepted bool) {
	opponent, _ := g.GetOpponentPlayer(p)
	if g.drawOffered != opponent {
		return
	}

	g.drawOffered = nil
	if !accepted {
		opponent.client.Send(map[string]any{
			"type": "respondDraw",
			"payload": map[string]string{
				"result": "rejected",
			},
		})
		return
	}
	g.Game.Exit()
	g.endGame(chess.Result{Reason: chess.Draw, WinnerColor: chess.Empty})
}

func (g *OnlineGame) Exit(p *onlinePlayer) error {
	color, err := g.getPlayerColor(p)
	if err != nil {
		return err
	}
	g.Game.Exit()

	return g.endGame(chess.Result{
		WinnerColor: color.OppositeColor(),
		Reason:      chess.Abandoned,
	})
}

func (g *OnlineGame) endGame(result chess.Result) error {
	for _, p := range g.Players {
		p.client.Send(map[string]any{
			"type": EndGameClientEvent,
			"payload": map[string]any{
				"winner": result.WinnerColor.String(),
				"score":  10,
				"reason": string(result.Reason),
			},
		})
		p.currentGame = nil
		p.status = StatusConnected
	}

	player1 := g.Players[chess.White]
	player2 := g.Players[chess.Black]

	if player1.user.IsAuthenticated() && player2.user.IsAuthenticated() {
		game := types.Game{
			Id: primitive.NewObjectID(),
			Players: []types.Player{
				{UserId: player1.user.GetId(), Color: chess.White.String()},
				{UserId: player2.user.GetId(), Color: chess.Black.String()},
			},
			Winner: result.WinnerColor.String(),
			Reason: string(result.Reason),
		}
		ctx := context.Background()
		return g.Storage.InsertGame(ctx, &game)
	}
	return nil
}

func (g *OnlineGame) GetOpponentPlayer(p *onlinePlayer) (*onlinePlayer, error) {
	color, err := g.getPlayerColor(p)
	if err != nil {
		return nil, err
	}
	return g.Players[color.OppositeColor()], nil
}
