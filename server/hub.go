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

type GameHandler interface {
	Start()

	Register(*player)
	UnRegister(*player)

	Play(player *player, gameId string, move game.Move)
	ExitGame(gameId string, player *player)

	AddToWaitList(p *player)
	RemoveFromWaitList(p *player)
}

type gameHandler struct {
	games   map[string]game.Chess
	players map[string]*player

	waitList WaitList

	waitListCh     chan *player
	exitWaitListCh chan *player
	registerCh     chan *player
	unregisterCh   chan *player
	exitGameCh     chan exitGameMsg
	playMoveCh     chan playMoveMsg
}

func NewGameHandler(wl WaitList) GameHandler {
	h := &gameHandler{
		games:   map[string]game.Chess{},
		players: map[string]*player{},

		waitList:       wl,
		waitListCh:     make(chan *player),
		exitWaitListCh: make(chan *player),
		exitGameCh:     make(chan exitGameMsg),
		playMoveCh:     make(chan playMoveMsg),
		registerCh:     make(chan *player),
		unregisterCh:   make(chan *player),
	}

	return h
}

func (h *gameHandler) Register(p *player) {
	h.registerCh <- p
}
func (h *gameHandler) UnRegister(p *player) {
	h.unregisterCh <- p
}
func (h *gameHandler) Play(player *player, gameId string, move game.Move) {
	h.playMoveCh <- playMoveMsg{player: player, gameId: gameId, move: move}
}
func (h *gameHandler) ExitGame(gameId string, p *player) {
	h.exitGameCh <- exitGameMsg{player: p, gameId: gameId}
}
func (h *gameHandler) AddToWaitList(p *player) {
	h.waitListCh <- p
}
func (h *gameHandler) RemoveFromWaitList(p *player) {
	h.exitWaitListCh <- p
}
func (h *gameHandler) Start() {
	for {
		select {
		case p := <-h.registerCh:
			h.handleRegister(p)
		case p := <-h.unregisterCh:
			h.handleUnregister(p)
		case pm := <-h.playMoveCh:
			h.handlePlayerMove(pm.player, pm.move, pm.gameId)
		case p := <-h.waitListCh:
			h.handleWait(p)
		case p := <-h.exitWaitListCh:
			h.handleExitWaitList(p)
		case msg := <-h.exitGameCh:
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
	g := game.NewChess(map[string]game.Color{
		p1.GetId(): game.White,
		p2.GetId(): game.Black,
	})
	h.games[gameId] = g

	p1.Join <- joinMsg{
		GameId:   gameId,
		Opponent: p2.GetName(),
		Color:    game.White.String(),
	}
	p2.Join <- joinMsg{
		GameId:   gameId,
		Opponent: p1.GetName(),
		Color:    game.Black.String(),
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
				Winner: winner.String(),
			}
		}
	}

	delete(h.games, gameId)
}
func (h *gameHandler) handleExitWaitList(p *player) {
	if err := h.waitList.Remove(p); err != nil {
		p.Err <- err
	}
}
