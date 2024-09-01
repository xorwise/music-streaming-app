package middleware

import (
	"log/slog"
	"net/http"
	"time"
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
		next.ServeHTTP(w, r)
		lm.Log.Info("request", slog.String("method", r.Method), slog.String("url", r.URL.String()), slog.Duration("duration", time.Since(startTime)))
	})
}
