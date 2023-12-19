package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIService struct {
	Storage       Storage
	Authenticator Authenticator
}

func NewAPIs(storage Storage, auth Authenticator) *APIService {
	NewValidator()
	return &APIService{
		Storage:       storage,
		Authenticator: auth,
	}
}

func (s *APIService) RegistrationAPI(c echo.Context) error {
	userReq := RegistrationRequest{}
	if err := c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}
	if err := userReq.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	user := NewUser(userReq.Email, userReq.Name, userReq.Password)
	if err := s.Storage.InsertUser(c.Request().Context(), user); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "created"})
}

func (s *APIService) AuthenticationAPI(c echo.Context) error {
	authReq := AuthenticationRequest{}
	if err := c.Bind(&authReq); err != nil {
		return err
	}
	if err := authReq.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	user, err := s.Storage.AuthenticateUser(c.Request().Context(), authReq.Email, authReq.Password)
	if err != nil {
		if errors.Is(err, ErrAuthentication) {
			return c.JSON(
				http.StatusUnauthorized,
				map[string]string{"message": err.Error()},
			)
		}
		return err
	}

	token, err := s.Authenticator.ObtainToken(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (s *APIService) UsersAPI(c echo.Context) error {
	users, err := s.Storage.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}

// func (s *APIService) getMyUserHandler(c echo.Context) error {
// 	return writeJSON(w, http.StatusOK, myUser)
// }

func (s *APIService) AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenStr := c.Request().Header.Get("Authorization")
		user, err := s.Authenticator.Authenticate(c.Request().Context(), tokenStr)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}
		c.Set("UserId", user)
		return next(c)
	}
}
