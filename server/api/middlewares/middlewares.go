package middlewares

import (
	"database/sql"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/richardpanda/composition/server/api/types"
)

func DB(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"message": "Authorization header is required."})
			return
		}

		tokenString := authHeader[len("Bearer "):]

		t, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return types.JWTSecret, nil
		})

		if err != nil || !t.Valid {
			c.AbortWithStatusJSON(400, gin.H{"message": "Invalid token."})
			return
		}

		c.Set("user", t.Claims.(jwt.MapClaims))
		c.Next()
	}
}
