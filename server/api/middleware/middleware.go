package middleware

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/richardpanda/composition/server/api/types"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Authorization header is required.", http.StatusBadRequest)
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Malformed authorization header.", http.StatusBadRequest)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return types.JWTSecret, nil
		})

		if err != nil || !t.Valid {
			http.Error(w, "Invalid token.", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", t.Claims.(jwt.MapClaims))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
