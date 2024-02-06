package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sina-am/chess/chess"
)

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
    Player *player
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
    Move chess.Move
}

type GameHandler interface {
	Start()

	Register(*player)
	UnRegister(*player)

	Play(player *player, gameId string, move chess.Move)
	ExitGame(gameId string, player *player)

	AddToWaitList(p *player)
	RemoveFromWaitList(p *player)
}

type gameHandler struct {
	games   map[string]chess.Chess
	players map[string]*player

	waitList WaitList

    eventCh        chan EventMsg 
}

func NewGameHandler(wl WaitList) GameHandler {
	h := &gameHandler{
		games:   map[string]chess.Chess{},
		players: map[string]*player{},
		waitList:       wl,
        eventCh: make(chan EventMsg),
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
            Move: move,
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

func (h *gameHandler) AddToWaitList(p *player) {
    msg := EventMsg{
        Type: JoinWaitListEventType,
        Body: JoinWaitListEventMsg{
            Player: p,
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
                h.handleWait(body.Player)
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

func (h *gameHandler) publishGameFinished(g chess.Chess) {
	color := g.GetWinner()
	for _, playerId := range g.GetPlayers() {
		if p, ok := h.players[playerId]; ok {
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

	if g.IsFinished() {
		h.publishGameFinished(g)
	}
}

func (h *gameHandler) startNewGame(p1, p2 *player) {
	gameId := uuid.New().String()

	g := chess.NewSession([]string{p1.GetId(), p2.GetId()})

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
