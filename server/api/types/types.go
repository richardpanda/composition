package types

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type ArticlePreview struct {
	Username  string    `json:"username"`
	Title     string    `json:"title"`
	ID        int       `json:"article_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ErrorResponseBody struct {
	Message string `json:"message"`
}

type GetArticlesBody struct {
	ArticlePreviews []ArticlePreview `json:"article_previews"`
}

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type PostArticlesBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type SigninBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SigninResponseBody struct {
	Token string `json:"token"`
}

type SignupBody struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type SignupResponseBody struct {
	Token string `json:"token"`
}
