package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/utils"
)

type APIService struct {
	Storage       Storage
	Authenticator Authenticator
	Renderer      utils.Renderer
}

func NewAPIs(storage Storage, auth Authenticator, renderer utils.Renderer) *APIService {
	NewValidator()
	return &APIService{
		Storage:       storage,
		Authenticator: auth,
		Renderer:      renderer,
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

func (s *APIService) AuthenticationPOST(c echo.Context) error {
	authReq := AuthenticationRequest{}
	if err := c.Bind(&authReq); err != nil {
		return err
	}

	if err := authReq.Validate(); err != nil {
		return s.Renderer.Render(c.Response().Writer, "login.html", map[string]string{
			"error": err.Error(),
		})
	}
	user, err := s.Storage.AuthenticateUser(c.Request().Context(), authReq.Email, authReq.Password)
	if err != nil {
		if errors.Is(err, ErrAuthentication) {
			return s.Renderer.Render(c.Response().Writer, "login.html", map[string]string{
				"error": err.Error(),
			})
		}
		return err
	}

	if err := s.Authenticator.Login(c, user); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/")
}

func (s *APIService) AuthenticationAPI(c echo.Context) error {
	if c.Request().Method == http.MethodPost {
		return s.AuthenticationPOST(c)
	}
	return s.Renderer.Render(c.Response().Writer, "login.html", nil)
}

func (s *APIService) UsersAPI(c echo.Context) error {
	users, err := s.Storage.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}
