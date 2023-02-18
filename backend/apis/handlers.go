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
	userReq := &types.RegistrationRequest{}
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
	authReq := &types.AuthenticationRequest{}
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

	myUser.Gender = updateReq.Gender
	myUser.Picture = updateReq.Picture
	myUser.Name = updateReq.Name

	s.Database.UpdateUser(ctx, myUser)
	return writeJSON(w, http.StatusOK, myUser)
}

func (s *APIServer) getMyUserHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	myUser := ctx.Value(UserIdContext).(*types.User)
	return writeJSON(w, http.StatusOK, myUser)
}

func (s *APIServer) startGameHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return writeJSON(w, http.StatusOK, map[string]string{"message": "request created"})
}
