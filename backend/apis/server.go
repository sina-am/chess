package apis

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/engine"
	"go.uber.org/zap"
)

type APIServer struct {
	Addr          string
	Logger        *zap.SugaredLogger
	Upgrader      *websocket.Upgrader
	Database      database.Database
	Authenticator Authenticator
	Game          engine.Game
}

type RequestID string

type apiFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) makeAPIHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(context.TODO(), RequestID("RequestID"), rand.Intn(10000))
		if err := f(ctx, w, r); err != nil {
			httpError := &HTTPError{}
			if errors.As(err, &httpError) {
				writeJSON(w, httpError.StatusCode, httpError)
			} else {
				s.Logger.Errorw(err.Error())
			}
		}
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.Use(s.LoggerMiddleware)
	router.HandleFunc("/", s.makeAPIHandler(s.indexHandler))
	router.HandleFunc("/ws", s.makeAPIHandler(s.wsHandler))
	router.HandleFunc("/auth", s.makeAPIHandler(s.authenticationHandler)).Methods(http.MethodPost)
	router.HandleFunc("/users", s.makeAPIHandler(s.insertUserHandler)).Methods(http.MethodPost)
	router.HandleFunc(
		"/users/me",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.getMyUserHandler)),
	)
	router.HandleFunc(
		"/users",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.getAllUsersHandler)),
	).Methods(http.MethodGet)

	router.HandleFunc(
		"/games",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.startGameHandler)),
	).Methods(http.MethodPost)

	s.Logger.Infow("server is running", "address", s.Addr)

	return http.ListenAndServe(s.Addr, router)
}
