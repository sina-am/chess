package game

import (
	"fmt"
	"time"

	"github.com/sina-am/chess/chess"
)

type ClientEventType string

const (
	StartedClientEvent ClientEventType = "started"
	EndGameClientEvent ClientEventType = "ended"
	PlayedClientEvent  ClientEventType = "played"
)

type OnlineGame struct {
	Players map[chess.Color]*onlinePlayer
	Game    chess.Chess
}

func NewOnlineGame(p1, p2 *onlinePlayer, duration time.Duration) *OnlineGame {
	game := &OnlineGame{
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

	p1.client.Send(map[string]any{
		"type": StartedClientEvent,
		"payload": map[string]any{
			"opponent": p2.user.GetName(),
			"tile":     chess.White.String(),
		},
	})
	p2.client.Send(map[string]any{
		"type": StartedClientEvent,
		"payload": map[string]any{
			"opponent": p1.user.GetName(),
			"tile":     chess.Black.String(),
		},
	})

	return game
}

func (g *OnlineGame) GetWinner() chess.Color {
	return g.Game.GetWinner()
}

func (g *OnlineGame) GetPlayers() []*onlinePlayer {
	players := []*onlinePlayer{}
	for color := range g.Players {
		players = append(players, g.Players[color])
	}

	return players
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

	for _, pl := range g.GetPlayers() {
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

	if g.Game.IsFinished() {
		color := g.GetWinner()

		g.broadcastEndMsg(map[string]any{
			"type": EndGameClientEvent,
			"payload": map[string]any{
				"winner": color.String(),
				"score":  10,
				"reason": "Won",
			},
		})
	}
	return nil
}

func (g *OnlineGame) Exit(p *onlinePlayer) error {
	color, err := g.getPlayerColor(p)
	if err != nil {
		return err
	}
	winner, err := g.Game.Exit(color)
	if err != nil {
		return err
	}

	g.broadcastEndMsg(map[string]any{
		"type": EndGameClientEvent,
		"payload": map[string]any{
			"winner": winner.String(),
			"score":  10,
			"reason": "abandoned",
		},
	})

	return nil
}

func (g *OnlineGame) broadcastEndMsg(msg map[string]any) {
	for _, p := range g.GetPlayers() {
		p.client.Send(msg)
		p.currentGame = nil
		p.status = StatusConnected
	}
}
