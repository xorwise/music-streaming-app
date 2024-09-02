package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xorwise/music-streaming-service/internal/domain"
)

type jwtMiddleware struct {
	Secret     string
	Repository domain.UserRepository
}

func NewJWTMiddleware(secret string, ur domain.UserRepository) *jwtMiddleware {
	return &jwtMiddleware{
		Secret:     secret,
		Repository: ur,
	}
}

func (j *jwtMiddleware) LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			authHeader = r.Header.Get("Sec-Websocket-Protocol")
			if authHeader == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(j.Secret), nil
		})

		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id := int64(claims["id"].(float64))
		user, err := j.Repository.GetByID(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
