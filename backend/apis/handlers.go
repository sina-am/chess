package apis

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/types"
)

func methodNotAllowedHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return &HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		Message:    "error",
		Details:    "method not allowd",
	}
}

func notFoundHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return &HTTPError{
		StatusCode: http.StatusNotFound,
		Message:    "error",
		Details:    "not found",
	}
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

func (s *APIServer) getModel(r *http.Request, model types.RequestModel) error {
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(model); err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid json",
			Details:    err.Error(),
		}
	}

	if err := model.Validate(); err != nil {
		return &HTTPError{
			StatusCode: http.StatusBadRequest,
			Message:    "invalid model",
			Details:    err.Error(),
		}
	}

	return nil
}

func (s *APIServer) insertUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	userReq := &types.RegisterUserRequest{}
	if err := s.getModel(r, userReq); err != nil {
		return err
	}

	user := types.NewUser(userReq.Email, userReq.Password)
	if err := s.Database.InsertUser(ctx, user); err != nil {
		return err
	}

	return writeJSON(w, http.StatusCreated, map[string]string{"message": "created"})
}
func (s *APIServer) authenticationHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	authReq := &types.AuthenticateUserRequest{}
	if err := s.getModel(r, authReq); err != nil {
		return err
	}

	user, err := s.Database.AuthenticateUser(ctx, authReq.Email, authReq.Password)
	if err != nil {
		if errors.Is(err, database.ErrAuthentication) {
			return &HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
				Details:    err.Error(),
			}
		}
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

func (s *APIServer) updateMyUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	myUser := ctx.Value(UserIdContext).(*types.User)
	updateReq := &types.UpdateUserRequest{}
	if err := s.getModel(r, updateReq); err != nil {
		return err
	}
	s.Database.UpdateUser()
	return writeJSON(w, http.StatusOK, user)
}

func (s *APIServer) getMyUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	myUser := ctx.Value(UserIdContext).(*types.User)
	return writeJSON(w, http.StatusOK, myUser)
}

func (s *APIServer) startGameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	gameReq := &types.StartGameRequest{}
	if err := s.getModel(r, gameReq); err != nil {
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
