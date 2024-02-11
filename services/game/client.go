package game

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/auth"
	"github.com/sina-am/chess/chess"
)

var (
	ErrInvalidType    = errors.New("invalid type")
	ErrInvalidPayload = errors.New("invalid payload")
)

type message struct {
	Type    string `json:"type"`
	Payload json.RawMessage
}

type Client interface {
	Send(msg any)
	SendErr(err error)
	Close()
}

type WSClient struct {
	conn *websocket.Conn
	user auth.User

	send  chan any
	err   chan error
	close chan error

	gameHandler GameHandler
	msgHandler  map[string]func(message) error
}

func NewWSClient(conn *websocket.Conn, gamHandler GameHandler, user auth.User) *WSClient {
	client := &WSClient{
		conn:        conn,
		user:        user,
		gameHandler: gamHandler,
		send:        make(chan any),
		err:         make(chan error),
		close:       make(chan error),
	}
	client.msgHandler = map[string]func(message) error{
		"start":       client.handleStart,
		"play":        client.handlePlay,
		"exit":        client.handleExit,
		"offerDraw":   client.handleOfferDraw,
		"respondDraw": client.handleRespondDraw,
	}
	return client
}

func (p *WSClient) Close() {
	close(p.send)
}

func (p *WSClient) Send(msg any) {
	p.send <- msg
}

func (p *WSClient) SendErr(err error) {
	p.err <- err
}

func (p *WSClient) StartLoop(ctx context.Context) {
	go p.readConn(ctx)
	go p.writeConn(ctx)
}

func (p *WSClient) readConn(ctx context.Context) {
	defer func() {
		log.Printf("player %s disconnected", p.conn.RemoteAddr())
		p.gameHandler.UnRegister(p)
		p.conn.Close()
	}()

	for {
		msg := message{}
		if err := p.conn.ReadJSON(&msg); err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				break
			}
			log.Printf("websocket error: %v", err)
			p.err <- err
			break
		}

		if err := p.handleMessage(msg); err != nil {
			p.err <- err
			continue
		}
	}
}

func (p *WSClient) writeConn(ctx context.Context) {
	defer func() {
		log.Printf("Go routine exited")
	}()

	for {
		select {
		case msg, ok := <-p.send:
			if !ok {
				return
			}
			if err := p.conn.WriteJSON(msg); err != nil {
				return
			}
		case err, ok := <-p.err:
			if !ok {
				break
			}
			if err := p.conn.WriteJSON(map[string]string{"error": err.Error()}); err != nil {
				return
			}
		case <-p.close:
			return
		}
	}
}

type StartGameMessage struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Duration int    `json:"duration"`
}

func (p *WSClient) handleMessage(msg message) error {
	handler, found := p.msgHandler[msg.Type]
	if !found {
		return ErrInvalidType
	}
	return handler(msg)
}

func (p *WSClient) handleStart(msg message) error {
	payload := StartGameMessage{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return ErrInvalidPayload
	}

	duration := 0 * time.Minute
	switch payload.Duration {
	case 10:
		duration = 10 * time.Minute
	case 5:
		duration = 5 * time.Minute
	case 3:
		duration = 3 * time.Minute
	case 1:
		duration = 1 * time.Minute
	default:
		return ErrInvalidPayload
	}

	p.gameHandler.AddToWaitList(p, GameSetting{Duration: duration})
	return nil
}

type PlayGameMessage struct {
	Move chess.Move `json:"move"`
}

func (p *WSClient) handlePlay(msg message) error {
	payload := PlayGameMessage{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	p.gameHandler.Play(p, payload.Move)
	return nil
}

func (p *WSClient) handleExit(msg message) error {
	p.gameHandler.Exit(p)
	return nil
}

func (p *WSClient) handleOfferDraw(msg message) error {
	p.gameHandler.OfferDraw(p)
	return nil
}

type respondDrawMessage struct {
	Result string
}

func (p *WSClient) handleRespondDraw(msg message) error {
	payload := respondDrawMessage{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return err
	}
	if payload.Result == "accepted" {
		p.gameHandler.RespondDraw(p, true)
	} else {
		p.gameHandler.RespondDraw(p, false)
	}

	return nil
}
