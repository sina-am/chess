package apis

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/engine"
	"github.com/sina-am/chess/types"
)

func (s *APIServer) startNewGameHandler(ctx context.Context, myUser *types.User, ws *websocket.Conn) (engine.OnlinePlayer, error) {
	startGameRequest := &types.StartGameRequest{}
	if err := ws.ReadJSON(startGameRequest); err != nil {
		return nil, &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid request",
			Details:    err.Error(),
		}
	}

	player, err := s.Game.StartGame(ctx, myUser)
	if err != nil {
		return nil, &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "can't start a new game",
			Details:    err.Error(),
		}
	}

	game := player.GetGame()
	err = ws.WriteJSON(game)
	return player, err
}

func (s *APIServer) playHandler(ctx context.Context, ws *websocket.Conn, player engine.OnlinePlayer) error {
	if err := player.WaitForMyTurn(ctx); err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "something went wrong",
			Details:    err.Error(),
		}
	}
	if err := ws.WriteJSON(map[string]string{"status": "your turn"}); err != nil {
		return err
	}

	playReq := &types.PlayRequest{}
	if err := ws.ReadJSON(playReq); err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "can't read json for play request",
			Details:    err.Error(),
		}
	}

	if err := playReq.Validate(); err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "can't validate play request",
			Details:    err.Error(),
		}
	}

	if err := player.Play(ctx, playReq.From, playReq.To); err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "can't play",
			Details:    err.Error(),
		}
	}

	return ws.WriteJSON(playReq)
}

func (s *APIServer) wsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newRequestContext(context.TODO())
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Logger.Errorw(err.Error())
		return
	}
	defer ws.Close()

	myUser, err := s.Authenticator.Authenticate(ctx, r.URL.Query().Get("token"))
	if err != nil {
		ws.WriteJSON(&HTTPError{
			StatusCode: http.StatusUnauthorized,
			Message:    "unauthorized",
			Details:    err.Error(),
		})
		s.Logger.Errorw(err.Error())
		return
	}

	for {
		player, err := s.startNewGameHandler(ctx, myUser, ws)
		if err != nil {
			if err := ws.WriteJSON(err); err != nil {
				return
			}
			continue
		}
		for {
			if err := s.playHandler(ctx, ws, player); err != nil {
				if err := ws.WriteJSON(err); err != nil {
					return
				}
				continue
			}
		}
	}
}
