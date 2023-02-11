package apis

import (
	"context"
	"net/http"
)

func (s *APIServer) CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		s.Logger.Infow("HELLO")
		w.Header().Set("Access-Control-Allow-Origin", req.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "content-type, authorization")

		if req.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, req)
	})
}

func (s *APIServer) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		next.ServeHTTP(w, req)
		s.Logger.Infow(
			"request",
			"method", req.Method,
			"url", req.URL,
			"proto", req.Proto,
			"status", w.Header().Get("status"),
		)
	})
}

var UserIdContext string = "user"

func (s *APIServer) AuthenticationMiddleware(f apiFunc) apiFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		user, err := s.Authenticator.Authenticate(ctx, r)
		if err != nil {
			return &HTTPError{
				StatusCode: http.StatusUnauthorized,
				Message:    "unauthorized",
				Details:    err.Error(),
			}
		}
		ctx = context.WithValue(ctx, UserIdContext, user)
		return f(ctx, w, r)
	}
}
