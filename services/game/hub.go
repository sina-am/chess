package game

import (
	"fmt"
	"log"
	"time"

	"github.com/sina-am/chess/auth"
	"github.com/sina-am/chess/chess"
)

type PlayerStatus int

const (
	StatusWaiting   PlayerStatus = 0
	StatusPlaying   PlayerStatus = 1
	StatusConnected PlayerStatus = 2
)

type onlinePlayer struct {
	client      Client
	user        auth.User
	status      PlayerStatus
	currentGame *OnlineGame
}

type onlinePlayerStorage struct {
	players map[Client]*onlinePlayer
}

func NewOnlinePlayerStorage() *onlinePlayerStorage {
	return &onlinePlayerStorage{
		players: map[Client]*onlinePlayer{},
	}
}

func (s *onlinePlayerStorage) Add(c Client, p *onlinePlayer) {
	if oldPlayer, ok := s.players[c]; ok {
		oldPlayer.client.Close()
	}
	s.players[c] = p
}

func (s *onlinePlayerStorage) Get(c Client) *onlinePlayer {
	return s.players[c]
}

func (s *onlinePlayerStorage) Remove(c Client) {
	p, ok := s.players[c]
	if !ok {
		return
	}

	delete(s.players, c)
	p.client.Close()
}

type EventType int

const (
	RegisterEventType EventType = iota
	UnRegisterEvent
	PlayEvent
	JoinWaitListEvent
	LeaveWaitListEvent
	ExitEvent
)

type EventMsg struct {
	Type EventType
	Body any
}

type RegisterEventMsg struct {
	Player Client
	User   auth.User
}

type UnRegisterEventMsg struct {
	Player Client
}

type JoinWaitListEventMsg struct {
	Player      Client
	GameSetting GameSetting
}

type LeaveWaitListEventMsg struct {
	Player Client
}

type ExitEventMsg struct {
	Player Client
}

type PlayEventMsg struct {
	Player Client
	Move   chess.Move
}

type GameSetting struct {
	Duration time.Duration
}
type GameHandler interface {
	Start()

	Register(p Client, user auth.User)
	UnRegister(Client)

	Play(client Client, move chess.Move)
	Exit(clien Client)

	AddToWaitList(p Client, gs GameSetting)
	RemoveFromWaitList(p Client)
}

type gameHandler struct {
	players  *onlinePlayerStorage
	waitList WaitList
	eventCh  chan EventMsg
}

func NewGameHandler(wl WaitList) GameHandler {
	h := &gameHandler{
		players:  NewOnlinePlayerStorage(),
		waitList: wl,
		eventCh:  make(chan EventMsg),
	}

	return h
}
func (h *gameHandler) Register(p Client, user auth.User) {
	msg := EventMsg{
		Type: RegisterEventType,
		Body: RegisterEventMsg{Player: p, User: user},
	}
	h.eventCh <- msg
}

func (h *gameHandler) UnRegister(p Client) {
	msg := EventMsg{
		Type: UnRegisterEvent,
		Body: UnRegisterEventMsg{Player: p},
	}
	h.eventCh <- msg
}

func (h *gameHandler) Play(p Client, move chess.Move) {
	msg := EventMsg{
		Type: PlayEvent,
		Body: PlayEventMsg{
			Player: p,
			Move:   move,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) Exit(p Client) {
	msg := EventMsg{
		Type: ExitEvent,
		Body: ExitEventMsg{
			Player: p,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) AddToWaitList(p Client, gs GameSetting) {
	msg := EventMsg{
		Type: JoinWaitListEvent,
		Body: JoinWaitListEventMsg{
			Player:      p,
			GameSetting: gs,
		},
	}
	h.eventCh <- msg
}

func (h *gameHandler) RemoveFromWaitList(p Client) {
	msg := EventMsg{
		Type: LeaveWaitListEvent,
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
			h.handleRegister(body.Player, body.User)
		case UnRegisterEvent:
			body := event.Body.(UnRegisterEventMsg)
			h.handleUnregister(body.Player)
		case PlayEvent:
			body := event.Body.(PlayEventMsg)
			h.handlePlayerMove(body.Player, body.Move)
		case JoinWaitListEvent:
			body := event.Body.(JoinWaitListEventMsg)
			h.handleWait(body.Player, body.GameSetting)
		case LeaveWaitListEvent:
			body := event.Body.(LeaveWaitListEventMsg)
			h.handleExitWaitList(body.Player)
		case ExitEvent:
			body := event.Body.(ExitEventMsg)
			h.handleExit(body.Player)
		}
	}
}
func (h *gameHandler) handleRegister(c Client, user auth.User) {
	op := &onlinePlayer{client: c, status: StatusConnected, user: user}
	h.players.Add(c, op)
}

func (h *gameHandler) handleUnregister(p Client) {
	h.handleExit(p)
	h.players.Remove(p)
}

func (h *gameHandler) handlePlayerMove(c Client, move chess.Move) {
	player := h.players.Get(c)
	if player == nil {
		log.Printf("player with client %v is not in the players list", c)
		return
	}

	if player.status != StatusPlaying {
		player.client.SendErr(fmt.Errorf("you're not in any game"))
		return
	}

	if err := player.currentGame.Play(player, move); err != nil {
		log.Printf("onlineGame.Play: %s", err.Error())
		return
	}
}

func createWaitListKey(gs GameSetting) string {
	return fmt.Sprintf("<%d>", gs.Duration)
}

func (h *gameHandler) handleWait(c Client, gs GameSetting) {
	player := h.players.Get(c)
	if player == nil {
		log.Printf("player with client %v is not in the players list", c)
		return
	}

	if player.status == StatusWaiting {
		c.SendErr(fmt.Errorf("already in a waiting list"))
		return
	}

	if player.status == StatusPlaying {
		c.SendErr(fmt.Errorf("already in a game"))
		return
	}

	c2, err := h.waitList.Pop(createWaitListKey(gs))

	// Wait, list is empty
	if err != nil {
		if err := h.waitList.Add(createWaitListKey(gs), c); err != nil {
			c.SendErr(err)
		}
		player.status = StatusWaiting
		return
	}

	player2 := h.players.Get(c2)
	NewOnlineGame(player, player2, gs.Duration)
}

func (h *gameHandler) handleExit(c Client) {
	player := h.players.Get(c)
	if player == nil {
		log.Printf("player with client %v is not in the players list", c)
		return
	}

	if player.status == StatusWaiting {
		h.waitList.Remove(c)
		player.status = StatusConnected
	} else if player.status == StatusPlaying {
		g := player.currentGame
		g.Exit(player)
	}
}
func (h *gameHandler) handleExitWaitList(c Client) {
	if err := h.waitList.Remove(c); err != nil {
		c.SendErr(err)
	}
}
