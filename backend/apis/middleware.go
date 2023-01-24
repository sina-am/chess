package apis

import (
	"context"
	"net/http"
)

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
