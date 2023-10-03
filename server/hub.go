package server

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/sina-am/chess/game"
)

type playMoveMsg struct {
	player *player
	gameId string
	move   game.Move
}

type exitGameMsg struct {
	gameId string
	player *player
}

type gameHandler struct {
	games   map[string]game.Chess
	players map[string]*player

	waitList   WaitList
	wait       chan *player
	exitQueue  chan *player
	register   chan *player
	unregister chan *player

	exitGame chan exitGameMsg
	playMove chan playMoveMsg
}

func NewGameHandler(wl WaitList) *gameHandler {
	h := &gameHandler{
		games:   map[string]game.Chess{},
		players: map[string]*player{},

		waitList:   wl,
		wait:       make(chan *player),
		exitQueue:  make(chan *player),
		exitGame:   make(chan exitGameMsg),
		playMove:   make(chan playMoveMsg),
		register:   make(chan *player),
		unregister: make(chan *player),
	}

	return h
}
func (h *gameHandler) Start() {
	for {
		select {
		case p := <-h.register:
			h.handleRegister(p)
		case p := <-h.unregister:
			h.handleUnregister(p)
		case pm := <-h.playMove:
			h.handlePlayerMove(pm.player, pm.move, pm.gameId)
		case p := <-h.wait:
			h.handleWait(p)
		case p := <-h.exitQueue:
			h.handleExitQueue(p)
		case msg := <-h.exitGame:
			h.handleExitGame(msg.player, msg.gameId)
		}
	}
}
func (h *gameHandler) handleUnregister(p *player) {
	playerId := p.GetId()
	if _, ok := h.players[playerId]; ok {
		h.waitList.FindAndDelete(p)
		for gameId, g := range h.games {
			if g.InGame(playerId) {
				h.handleExitGame(p, gameId)
			}
		}
		delete(h.players, playerId)
		close(p.Send)
		close(p.Err)
		close(p.Join)
		close(p.End)
	}
}
func (h *gameHandler) handleRegister(p *player) {
	h.players[p.GetId()] = p
}

func (h *gameHandler) handlePlayerMove(p *player, move game.Move, gameId string) {
	g, found := h.games[gameId]
	if !found {
		p.Err <- fmt.Errorf("game not found")
		return
	}

	if err := g.Play(p.GetId(), move); err != nil {
		if errors.Is(err, game.ErrGameEnd) {
			t, _ := g.GetWinner()
			for _, playerId := range g.GetPlayers() {
				if p, ok := h.players[playerId]; ok {
					p.End <- endMsg{
						Winner: t,
						Score:  10,
						Reason: "Won",
					}
				}
			}
		} else {
			p.Err <- err
		}
		return
	}

	p.Send <- map[string]string{
		"message": fmt.Sprintf("you played %d", move),
	}
	for _, playerId := range g.GetPlayers() {
		if playerId != p.GetId() {
			if p, ok := h.players[playerId]; ok {
				p.Send <- map[string]any{
					"type": "played",
					"payload": map[string]any{
						"player": p.GetName(),
						"move":   move,
					},
				}
			}
		}
	}
}

func (h *gameHandler) startNewGame(p1, p2 *player) {
	gameId := uuid.New().String()
	g := game.NewChess([]string{p1.GetId(), p2.GetId()})
	h.games[gameId] = g

	p1.Join <- joinMsg{
		GameId:   gameId,
		Opponent: p2.GetName(),
		Color:    g.GetPlayerColor(p1.GetId()).String(),
	}
	p2.Join <- joinMsg{
		GameId:   gameId,
		Opponent: p1.GetName(),
		Color:    g.GetPlayerColor(p2.GetId()).String(),
	}
}

func (h *gameHandler) handleWait(player1 *player) {
	player2, err := h.waitList.Pop()

	// Wait list is empty
	if err != nil {
		if err := h.waitList.Add(player1); err != nil {
			player1.Err <- err
		}
		return
	}

	h.startNewGame(player1, player2)
}
func (h *gameHandler) handleExitGame(p *player, gameId string) {
	g := h.games[gameId]

	winner, err := g.Exit(p.GetId())
	if err != nil {
		return
	}

	for _, playerId := range g.GetPlayers() {
		if p, ok := h.players[playerId]; ok {
			p.End <- endMsg{
				Reason: "abundant",
				Score:  10,
				Winner: winner,
			}
		}
	}

	delete(h.games, gameId)
}
func (h *gameHandler) handleExitQueue(p *player) {
	if err := h.waitList.Remove(p); err != nil {
		p.Err <- err
	}
}
