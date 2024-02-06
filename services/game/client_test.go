package game

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type errorMsg struct {
	Err string `json:"error"`
}

func TestSendInvalidText(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}
	defer ws.Close()

	ws.WriteMessage(websocket.TextMessage, []byte("invalid data"))

	_, _, err = ws.ReadMessage()

	if err == nil {
		t.Error("expected to be closed")
	}
}
func TestInvalidMessageType(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}
	defer ws.Close()

	ws.WriteJSON(map[string]string{"type": "invalid"})

	msg := errorMsg{}
	if err := ws.ReadJSON(&msg); err != nil {
		assert.Equal(t, ErrInvalidType.Error(), msg.Err)
	}
}
func TestInvalidMessagePayload(t *testing.T) {
	ws, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		t.Error(err)
	}
	defer ws.Close()

	ws.WriteJSON(map[string]string{"type": "start"})

	msg := errorMsg{}
	if err := ws.ReadJSON(&msg); err != nil {
		assert.Equal(t, ErrInvalidPayload.Error(), msg.Err)
	}
}
func TestStartGame(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		dialer := websocket.DefaultDialer
		dialer.HandshakeTimeout = 0
		ws, _, err := dialer.Dial("ws://localhost:8080/ws", nil)
		ws.SetReadDeadline(time.Time{})
		if err != nil {
			t.Error(err)
		}
		defer ws.Close()

		ws.WriteJSON(map[string]any{
			"type": "start",
			"payload": map[string]any{
				"name":     "anonymous",
				"duration": 10,
			},
		})

		_, data, err := ws.ReadMessage()
		if err != nil {
			return
		}
		fmt.Printf("%s\n", data)
	}()
	go func() {
		defer wg.Done()
		dialer := websocket.DefaultDialer
		dialer.HandshakeTimeout = 0
		ws, _, err := dialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			t.Error(err)
		}
		defer ws.Close()

		ws.WriteJSON(map[string]any{
			"type": "start",
			"payload": map[string]any{
				"name":     "anonymous",
				"duration": 10,
			},
		})

		_, data, err := ws.ReadMessage()
		if err != nil {
			return
		}
		fmt.Printf("%s\n", data)
	}()

	wg.Wait()
}
