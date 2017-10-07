package types

import (
	"os"

	jwt "github.com/dgrijalva/jwt-go"
)

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type SignupBody struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type SignupResponse struct {
	Token string `json:"token"`
}