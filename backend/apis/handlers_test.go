package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/service"
	"github.com/sina-am/chess/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func getTestAPIServer() APIServer {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewMemoryDatabase(context.TODO())
	gameSrv, err := service.NewGameService(db)
	if err != nil {
		log.Fatal(err)
	}
	types.NewValidator()
	return APIServer{
		Addr:          ":8080",
		Game:          gameSrv,
		Logger:        logger.Sugar(),
		Database:      db,
		Authenticator: NewJWTAuthentication("verysecretkey", db),
		Upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}
func newRequest(method string, url string, body any) *http.Request {
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, url, bytes.NewReader(jsonBody))
	return req
}

func TestAuthenticateUserHandler(t *testing.T) {
	server := getTestAPIServer()

	user := types.NewUser("test@gmail.com", "test")
	server.Database.InsertUser(context.TODO(), user)

	t.Run("status code", func(t *testing.T) {
		request := newRequest(
			http.MethodPost,
			"/auth",
			types.AuthenticationRequest{
				Email:    "test@gmail.com",
				Password: "test",
			},
		)
		response := httptest.NewRecorder()

		err := server.authenticationHandler(context.TODO(), response, request)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})
	t.Run("checking token", func(t *testing.T) {
		request := newRequest(
			http.MethodPost,
			"/auth",
			types.AuthenticationRequest{
				Email:    "test@gmail.com",
				Password: "test",
			},
		)
		response := httptest.NewRecorder()

		err := server.authenticationHandler(context.TODO(), response, request)
		assert.Nil(t, err)

		tokenRes := map[string]string{}
		assert.Nil(t, json.NewDecoder(response.Body).Decode(&tokenRes))

		userRes, err := server.Authenticator.Authenticate(context.TODO(), tokenRes["token"])
		assert.Nil(t, err)
		assert.Equal(t, user, userRes)
	})

	t.Run("invalid email", func(t *testing.T) {
		request := newRequest(
			http.MethodPost,
			"/auth",
			types.AuthenticationRequest{
				Email:    "wrongemail@gmail.com",
				Password: "test",
			},
		)
		response := httptest.NewRecorder()

		err := server.authenticationHandler(context.TODO(), response, request)
		httpErr := &HTTPError{}
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusUnauthorized, httpErr.StatusCode)
	})
	t.Run("invalid password", func(t *testing.T) {
		request := newRequest(
			http.MethodPost,
			"/auth",
			types.AuthenticationRequest{
				Email:    "test@gmail.com",
				Password: "wrongpassword",
			},
		)
		response := httptest.NewRecorder()

		err := server.authenticationHandler(context.TODO(), response, request)
		httpErr := &HTTPError{}
		assert.ErrorAs(t, err, &httpErr)
		assert.Equal(t, http.StatusUnauthorized, httpErr.StatusCode)
	})
}

func TestGetAllUsersHandler(t *testing.T) {
	server := getTestAPIServer()

	users := []*types.User{
		types.NewUser("test1@gmail.com", "test"),
		types.NewUser("test2@gmail.com", "test"),
	}
	server.Database.InsertUser(context.TODO(), users[0])
	server.Database.InsertUser(context.TODO(), users[1])

	t.Run("returns all users in database", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/users", nil)
		response := httptest.NewRecorder()

		err := server.getAllUsersHandler(context.TODO(), response, request)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, "application/json", response.Header().Get("content-type"))
		responseUsers := []*types.User{}
		assert.Nil(t, json.NewDecoder(response.Body).Decode(&responseUsers))
	})
}
