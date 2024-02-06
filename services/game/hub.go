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

type EventType int

const (
	RegisterEventType EventType = iota
	UnRegisterEventType
	PlayEventType
	JoinWaitListEventType
	LeaveWaitListEventType
	ExitGameEventType
)

type EventMsg struct {
	Type EventType
	Body any
}

type RegisterEventMsg struct {
	Player *player
}

type UnRegisterEventMsg struct {
	Player *player
}

type JoinWaitListEventMsg struct {
	Player      *player
	GameSetting GameSetting
}

type LeaveWaitListEventMsg struct {
	Player *player
}

type ExitGameEventMsg struct {
	Player *player
	GameId string
}

type PlayEventMsg struct {
	Player *player
	GameId string
	Move   chess.Move
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

	eventCh chan EventMsg
}

func NewGameHandler(wl WaitList) GameHandler {
	h := &gameHandler{
		players:  NewOnlinePlayerStorage(),
		games:    map[string]chess.Chess{},
		waitList: wl,
		eventCh:  make(chan EventMsg),
	}

	return h
}
func (h *gameHandler) Register(p *player) {
	msg := EventMsg{
		Type: RegisterEventType,
		Body: RegisterEventMsg{Player: p},
	}
	h.eventCh <- msg
}

func (h *gameHandler) UnRegister(p *player) {
	msg := EventMsg{
		Type: UnRegisterEventType,
		Body: UnRegisterEventMsg{Player: p},
	}
	h.eventCh <- msg
}

func (h *gameHandler) Play(p *player, gameId string, move chess.Move) {
	msg := EventMsg{
		Type: PlayEventType,
		Body: PlayEventMsg{
			Player: p,
			GameId: gameId,
			Move:   move,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) ExitGame(gameId string, p *player) {
	msg := EventMsg{
		Type: ExitGameEventType,
		Body: ExitGameEventMsg{
			Player: p,
			GameId: gameId,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) AddToWaitList(p *player, gs GameSetting) {
	msg := EventMsg{
		Type: JoinWaitListEventType,
		Body: JoinWaitListEventMsg{
			Player:      p,
			GameSetting: gs,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) RemoveFromWaitList(p *player) {
	msg := EventMsg{
		Type: LeaveWaitListEventType,
		Body: LeaveWaitListEventMsg{
			Player: p,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) Start() {
	for {
		event := <-h.eventCh
		switch event.Type {
		case RegisterEventType:
			body := event.Body.(RegisterEventMsg)
			h.handleRegister(body.Player)
		case UnRegisterEventType:
			body := event.Body.(UnRegisterEventMsg)
			h.handleUnregister(body.Player)
		case PlayEventType:
			body := event.Body.(PlayEventMsg)
			h.handlePlayerMove(body.Player, body.Move, body.GameId)
		case JoinWaitListEventType:
			body := event.Body.(JoinWaitListEventMsg)
			h.handleWait(body.Player, body.GameSetting)
		case LeaveWaitListEventType:
			body := event.Body.(LeaveWaitListEventMsg)
			h.handleExitWaitList(body.Player)
		case ExitGameEventType:
			body := event.Body.(ExitGameEventMsg)
			h.handleExitGame(body.Player, body.GameId)
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
