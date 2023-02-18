package apis

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
)

func newRequestContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, RequestID("RequestID"), rand.Intn(10000))
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) makeAPIHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := newRequestContext(context.TODO())
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
