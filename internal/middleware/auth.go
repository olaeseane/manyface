package middleware

import (
	"context"
	"net/http"

	"go.uber.org/zap"
	"manyface.net/internal/session"
)

var (
	noAuthURLs = map[string]struct{}{
		"/api/v1beta1/reg":   struct{}{},
		"/api/v1beta1/login": struct{}{},
		"/api/v2beta1/user":  struct{}{},
	}
)

func Auth(logger *zap.SugaredLogger, sm *session.SessionManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := noAuthURLs[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		sessID := r.Header.Get("session-id")
		sess, err := sm.Check(sessID)
		if err != nil {
			http.Error(w, "Authentication error", http.StatusUnauthorized)
			logger.Errorf("Session %v not found", sessID)
			return
		}
		ctx := context.WithValue(r.Context(), session.SessKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
