package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/sina-am/chess/chess"
	"github.com/sina-am/chess/types"
	"nhooyr.io/websocket"
)

type message struct {
	Type    types.ClientEventType `json:"type"`
	Payload json.RawMessage
}

type OnlineChessClient struct {
	ws     *websocket.Conn
	engine *chess.ChessEngine
	ui     *ChessUI

	me       types.Player
	opponent types.Player
}

func NewOnlineChessClient(ws *websocket.Conn) *OnlineChessClient {
	return &OnlineChessClient{
		ws:     ws,
		engine: chess.NewEngine(),
	}
}

func (game *OnlineChessClient) sendStartEvent(ctx context.Context) (types.StartGameMsgOut, error) {
	startMsg, _ := json.Marshal(map[string]any{
		"type": types.StartServerEvent,
		"payload": types.StartGameMsgIn{
			Duration: 10,
		},
	})
	game.ws.Write(ctx, websocket.MessageText, startMsg)
	_, msgBytes, err := game.ws.Read(ctx)
	if err != nil {
		return types.StartGameMsgOut{}, fmt.Errorf("websocket read error: %s", err)
	}
	msg := types.StartGameMsgOut{}
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		return types.StartGameMsgOut{}, fmt.Errorf("json unmarshal error: %s", err)
	}

	return msg, nil
}

func (game *OnlineChessClient) eventListener(ctx context.Context, wg *sync.WaitGroup) {
	for {
		_, msgBytes, err := game.ws.Read(ctx)
		if err != nil {
			fmt.Println(err)
			break
		}

		msg := message{}
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			fmt.Println(err)
			break
		}

		switch msg.Type {
		case types.PlayedClientEvent:
			payload := types.PlayGamePayloadMsgOut{}
			json.Unmarshal(msg.Payload, &payload)
			if err := game.engine.Play(game.opponent.Color, payload.Move); err != nil {
				fmt.Println(err)
				return
			}
			game.ui.Render()
			break
		case types.EndGameClientEvent:
			payload := types.EndGamePayloadMsgOut{}
			json.Unmarshal(msg.Payload, &payload)
			fmt.Println(payload.Reason, payload.Winner)
			game.ui.Finish(chess.Result{
				Reason:      payload.Reason,
				WinnerColor: payload.Winner,
			})
			return

		}
	}
	wg.Done()
}

func (game *OnlineChessClient) Start(ctx context.Context) error {
	msg, err := game.sendStartEvent(ctx)
	if err != nil {
		return err
	}
	game.me = msg.Payload.You
	game.opponent = msg.Payload.Opponent

	game.ui = NewChessUI(game.engine, game.me.Color)
	game.ui.HookPickupHandler(func(piece *chess.Piece) error {
		if piece.Color != game.me.Color {
			return fmt.Errorf("not your piece")
		}
		if game.me.Color != game.engine.GetTurn() {
			return fmt.Errorf("not your turn")
		}
		return nil
	})

	game.ui.HookDropHandler(func(piece *chess.Piece, x, y int) error {
		msg, _ := json.Marshal(map[string]any{
			"type": types.PlayServerEvent,
			"payload": types.PlayGameMsgIn{
				Move: chess.Move{
					From: piece.Location,
					To:   chess.Location{Row: y, Col: x},
				},
			},
		})
		if err := game.ws.Write(ctx, websocket.MessageText, msg); err != nil {
			return err
		}
		return nil
	})
	game.ui.Render()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go game.eventListener(ctx, wg)
	wg.Wait()

	return nil
}

func startOnlineGame(ctx context.Context) {
	ws, _, err := websocket.Dial(ctx, "ws://localhost:8080/ws", nil)
	if err != nil {
		fmt.Println("websocket error", err)
		return
	}
	defer ws.Close(websocket.StatusGoingAway, "BYE")

	game := NewOnlineChessClient(ws)
	game.Start(ctx)
}

func startOfflineGame(ctx context.Context) {
	engine := chess.NewEngine()
	ui := NewChessUI(engine, chess.White)
	ui.Render()
}

func main() {
	fmt.Println("WASM Go Initialized")
	var done chan struct{}
	startOnlineGame(context.Background())
	<-done
}
