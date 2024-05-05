package users

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sina-am/chess/core"
	"github.com/sina-am/chess/services/auth"
	"github.com/sina-am/chess/storage"
	"github.com/sina-am/chess/types"
)

type APIService struct {
	Storage       storage.Storage
	Authenticator auth.Authenticator
	Renderer      core.Renderer
}

func NewAPIService(storage storage.Storage, auth auth.Authenticator, renderer core.Renderer) *APIService {
	return &APIService{
		Authenticator: auth,
		Storage:       storage,
		Renderer:      renderer,
	}
}

func (s *APIService) RegistrationAPI(c echo.Context) error {
	userReq := RegistrationRequest{}
	if err := c.Bind(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}
	if err := c.Validate(&userReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	user := types.NewUser(userReq.Email, userReq.Name, userReq.Password)
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

	if err := c.Validate(&authReq); err != nil {
		return s.Renderer.Render(c, "login.html", map[string]any{
			"error": err.Error(),
		})
	}
	user, err := s.Storage.AuthenticateUser(c.Request().Context(), authReq.Email, authReq.Password)
	if err != nil {
		if errors.Is(err, storage.ErrAuthentication) {
			return s.Renderer.Render(c, "login.html", map[string]any{
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

func (s *APIService) AuthenticationGET(c echo.Context) error {
	return s.Renderer.Render(c, "login.html", nil)
}

func (s *APIService) UsersAPI(c echo.Context) error {
	users, err := s.Storage.GetAllUsers(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, users)
}
