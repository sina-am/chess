package apis

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/service"
	"go.uber.org/zap"
)

type APIServer struct {
	Addr          string
	Logger        *zap.SugaredLogger
	Upgrader      *websocket.Upgrader
	Database      database.Database
	Authenticator Authenticator
	Game          service.GameService
}

type RequestID string

type apiFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

func (s *APIServer) newRouter() *mux.Router {
	router := mux.NewRouter()
	router.MethodNotAllowedHandler = s.makeAPIHandler(methodNotAllowedHandler)
	router.NotFoundHandler = s.makeAPIHandler(notFoundHandler)
	router.Use(s.LoggerMiddleware)
	router.Use(s.CorsMiddleware)
	router.Use(mux.CORSMethodMiddleware(router))
	return router
}

func (s *APIServer) attachHandlers(router *mux.Router) *mux.Router {
	router.HandleFunc("/ws", s.wsHandler)
	router.HandleFunc("/auth", s.makeAPIHandler(s.authenticationHandler)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/users", s.makeAPIHandler(s.insertUserHandler)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc(
		"/users/me",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.getMyUserHandler)),
	).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc(
		"/users",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.getAllUsersHandler)),
	).Methods(http.MethodGet, http.MethodOptions)

	router.HandleFunc(
		"/games",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.startGameHandler)),
	).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc(
		"/users",
		s.makeAPIHandler(s.AuthenticationMiddleware(s.updateMyUserHandler)),
	).Methods(http.MethodPut, http.MethodOptions)
	return router
}

func (s *APIServer) Run() error {
	router := s.attachHandlers(s.newRouter())
	s.Logger.Infow("server is running", "address", s.Addr)
	return http.ListenAndServe(s.Addr, router)
}
