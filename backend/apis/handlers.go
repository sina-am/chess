package apis

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/types"
)

func (s *APIServer) indexHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	file, err := os.Open("./apis/index.html")
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "text/html")
	_, err = io.Copy(w, file)

	return err
}

func (s *APIServer) wsHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ws, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	request := &types.MovePieceRequest{}
	if err := ws.ReadJSON(request); err != nil {
		return err
	}

	return ws.WriteJSON(request)
}

func (s *APIServer) insertUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	userReq := &types.RegisterUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(userReq); err != nil {
		return err
	}

	if err := userReq.Validate(); err != nil {
		return err
	}

	user := types.NewUser(userReq.Email, userReq.Password)
	if err := s.Database.InsertUser(ctx, user); err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{"message": "created"})
}
func (s *APIServer) authenticationHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	authReq := &types.AuthenticateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(authReq); err != nil {
		return err
	}

	if err := authReq.Validate(); err != nil {
		return err
	}

	user, err := s.Database.GetUserByEmail(ctx, authReq.Email)
	if err != nil {
		return err
	}

	if err := core.VerifyPassword(authReq.Password, user.Password); err != nil {
		return err
	}

	token, err := s.Authenticator.ObtainToken(user)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (s *APIServer) getAllUsersHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	users, err := s.Database.GetAllUsers(ctx)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, users)
}

func (s *APIServer) getMyUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	myUser := ctx.Value(UserIdContext).(*types.User)
	return writeJSON(w, http.StatusOK, myUser)
}

func (s *APIServer) startGameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	gameReq := &types.StartGameRequest{}
	if err := json.NewDecoder(r.Body).Decode(gameReq); err != nil {
		return err
	}

	if err := gameReq.Validate(); err != nil {
		return err
	}

	to, err := s.Database.GetUserById(ctx, gameReq.PlayerUserId)
	if err != nil {
		return err
	}

	if err := s.Game.Request(ctx, ctx.Value(UserIdContext).(*types.User), to); err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]string{"message": "request created"})
}
