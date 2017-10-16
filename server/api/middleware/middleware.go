package middleware

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/richardpanda/composition/server/api/types"
	"github.com/richardpanda/composition/server/api/utils"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			utils.SetErrorResponse(w, 401, "Authorization header is required.")
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utils.SetErrorResponse(w, 400, "Malformed authorization header.")
			return
		}

		tokenString := authHeader[len("Bearer "):]

		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return types.JWTSecret, nil
		})

		if err != nil || !t.Valid {
			utils.SetErrorResponse(w, 400, "Invalid token.")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", t.Claims.(jwt.MapClaims))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
