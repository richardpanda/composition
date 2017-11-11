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

type GetArticleResponseBody struct {
	ID        int       `json:"article_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type GetArticlesResponseBody struct {
	ArticlePreviews []ArticlePreview `json:"article_previews"`
}

var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

type JWTClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type PostArticlesRequestBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type PostArticlesResponseBody struct {
	ArticleID int    `json:"article_id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
}

type SigninRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SigninResponseBody struct {
	Token string `json:"token"`
}

type SignupRequestBody struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type SignupResponseBody struct {
	Token string `json:"token"`
}
