package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/xorwise/music-streaming-service/internal/domain"
)

type loggingMiddleware struct {
	Log *slog.Logger
}

func NewLoggingMiddleware(log *slog.Logger) *loggingMiddleware {
	return &loggingMiddleware{
		Log: log,
	}
}

func (lm *loggingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		var user *domain.User
		val := r.Context().Value("user")
		if val != nil {
			user = val.(*domain.User)
		}

		next.ServeHTTP(w, r)
		if user != nil {
			lm.Log.Info("request", slog.String("method", r.Method), slog.String("url", r.URL.String()), slog.String("user", user.Username), slog.Duration("duration", time.Since(startTime)))
		} else {
			lm.Log.Info("request", slog.String("method", r.Method), slog.String("url", r.URL.String()), slog.Duration("duration", time.Since(startTime)))
		}
	})
}
