package middleware

import (
	"net/http"

	"golang.org/x/net/websocket"
)

type wsMiddleware struct{}

func NewWSMiddleware() *wsMiddleware {
	return &wsMiddleware{}
}

func (wm *wsMiddleware) Handle(next websocket.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := websocket.Server{
			Handler: websocket.Handler(next),
		}
		s.ServeHTTP(w, r)
	})
}
