package game

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sina-am/chess/chess"
)

type onlinePlayerStorage struct {
	mu      sync.Mutex
	players map[string]*player
}

func NewOnlinePlayerStorage() *onlinePlayerStorage {
	return &onlinePlayerStorage{
		mu:      sync.Mutex{},
		players: map[string]*player{},
	}
}

func (s *onlinePlayerStorage) Exists(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.players[id]
	return ok
}

func (s *onlinePlayerStorage) Add(id string, p *player) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if oldPlayer, ok := s.players[id]; ok {
		oldPlayer.Close <- fmt.Errorf("you connected with another connection")
	}
	s.players[id] = p
}

func (s *onlinePlayerStorage) Get(id string) *player {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.players[id]
}

func (s *onlinePlayerStorage) Remove(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p, ok := s.players[id]
	if !ok {
		return
	}

	delete(s.players, id)
	close(p.Send)
	close(p.Err)
	close(p.Join)
	close(p.End)
}

type playMoveMsg struct {
	player *player
	gameId string
	move   chess.Move
}
type startGameMsg struct {
	player      *player
	gameSetting GameSetting
}

type exitGameMsg struct {
	gameId string
	player *player
}

type GameSetting struct {
	Duration time.Duration
}
type GameHandler interface {
	Start()

	Register(*player)
	UnRegister(*player)

	Play(player *player, gameId string, move chess.Move)
	ExitGame(gameId string, player *player)

	AddToWaitList(p *player, gs GameSetting)
	RemoveFromWaitList(p *player)
}

type gameHandler struct {
	games   map[string]chess.Chess
	players *onlinePlayerStorage

	waitList WaitList

	waitListCh     chan startGameMsg
	exitWaitListCh chan *player
	registerCh     chan *player
	unregisterCh   chan *player
	exitGameCh     chan exitGameMsg
	playMoveCh     chan playMoveMsg
}

func NewGameHandler(wl WaitList) GameHandler {
	h := &gameHandler{
		games:   map[string]chess.Chess{},
		players: NewOnlinePlayerStorage(),

		waitList:       wl,
		waitListCh:     make(chan startGameMsg),
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
func (h *gameHandler) Play(player *player, gameId string, move chess.Move) {
	h.playMoveCh <- playMoveMsg{player: player, gameId: gameId, move: move}
}
func (h *gameHandler) ExitGame(gameId string, p *player) {
	h.exitGameCh <- exitGameMsg{player: p, gameId: gameId}
}
func (h *gameHandler) AddToWaitList(p *player, gs GameSetting) {
	h.waitListCh <- startGameMsg{player: p, gameSetting: gs}
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
		case msg := <-h.waitListCh:
			h.handleWait(msg.player, msg.gameSetting)
		case p := <-h.exitWaitListCh:
			h.handleExitWaitList(p)
		case msg := <-h.exitGameCh:
			h.handleExitGame(msg.player, msg.gameId)
		}
	}
}

func (h *gameHandler) handleUnregister(p *player) {
	playerId := p.GetId()
	if h.players.Exists(playerId) {
		h.waitList.FindAndDelete(p)
		for gameId, g := range h.games {
			if g.InGame(playerId) {
				h.handleExitGame(p, gameId)
			}
		}
		h.players.Remove(playerId)
	}
}

func (h *gameHandler) handleRegister(p *player) {
	h.players.Add(p.GetId(), p)
}

func (h *gameHandler) publishGameFinished(g chess.Chess) {
	color := g.GetWinner()
	for _, playerId := range g.GetPlayers() {
		if p := h.players.Get(playerId); p != nil {
			p.End <- endMsg{
				Winner: color.String(),
				Score:  10,
				Reason: "Won",
			}
		}
	}
}

func (h *gameHandler) handlePlayerMove(p *player, move chess.Move, gameId string) {
	g, found := h.games[gameId]
	if !found {
		p.Err <- fmt.Errorf("game not found")
		return
	}

	if err := g.Play(p.GetId(), move); err != nil {
		p.Err <- err
		return
	}

	p.Send <- map[string]string{
		"message": fmt.Sprintf("you played %d", move),
	}
	for _, playerId := range g.GetPlayers() {
		if playerId != p.GetId() {
			if p := h.players.Get(playerId); p != nil {
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

	if g.IsFinished() {
		h.publishGameFinished(g)
	}
}

func (h *gameHandler) startNewGame(p1, p2 *player, gs GameSetting) {
	gameId := uuid.New().String()

	g := chess.NewSession([]string{p1.GetId(), p2.GetId()}, gs.Duration)

	h.games[gameId] = g

	p1.Join <- joinMsg{
		GameId:   gameId,
		Opponent: p2.GetName(),
		Color:    chess.White.String(),
	}
	p2.Join <- joinMsg{
		GameId:   gameId,
		Opponent: p1.GetName(),
		Color:    chess.Black.String(),
	}
}

func createWaitListKey(gs GameSetting) string {
	return fmt.Sprintf("<%d>", gs.Duration)
}
func (h *gameHandler) handleWait(p *player, gs GameSetting) {
	p2, err := h.waitList.Pop(createWaitListKey(gs))

	// Wait list is empty
	if err != nil {
		if err := h.waitList.Add(createWaitListKey(gs), p); err != nil {
			p.Err <- err
		}
		return
	}

	h.startNewGame(p, p2, gs)
}
func (h *gameHandler) handleExitGame(p *player, gameId string) {
	g := h.games[gameId]

	winner, err := g.Exit(p.GetId())
	if err != nil {
		return
	}

	for _, playerId := range g.GetPlayers() {
		if p := h.players.Get(playerId); p != nil {
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
