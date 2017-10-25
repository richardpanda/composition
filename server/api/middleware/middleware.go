package middleware

import (
	"context"
	"net/http"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/richardpanda/composition/server/api/types"
	"github.com/richardpanda/composition/server/api/utils"
)

func IsAuthenticated(w http.ResponseWriter, r *http.Request) *http.Request {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		utils.SetErrorResponse(w, 401, "Authorization header is required.")
		return nil
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		utils.SetErrorResponse(w, 400, "Malformed authorization header.")
		return nil
	}

	tokenString := authHeader[len("Bearer "):]

	t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return types.JWTSecret, nil
	})

	if err != nil || !t.Valid {
		utils.SetErrorResponse(w, 400, "Invalid token.")
		return nil
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, "user", t.Claims.(jwt.MapClaims))
	return r.WithContext(ctx)
}
