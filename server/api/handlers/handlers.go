package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/richardpanda/composition/server/api/middleware"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
	"github.com/richardpanda/composition/server/api/utils"
	"golang.org/x/crypto/bcrypt"
)

func HandleArticles(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var err error
			page := 1
			pageNumStr := r.FormValue("page")

			if pageNumStr != "" {
				page, err = strconv.Atoi(pageNumStr)

				if err != nil {
					utils.SetErrorResponse(w, 500, err.Error())
					return
				}
			}

			rows, err := models.GetLatestArticlePreviews(db, page)

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			defer rows.Close()

			var articlePreviews []types.ArticlePreview

			for rows.Next() {
				var username, title string
				var id int
				var createdAt time.Time

				if err := rows.Scan(&username, &title, &id, &createdAt); err != nil {
					utils.SetErrorResponse(w, 500, err.Error())
					return
				}

				articlePreviews = append(articlePreviews, types.ArticlePreview{username, title, id, createdAt})
			}

			b, err := json.Marshal(types.GetArticlesBody{ArticlePreviews: articlePreviews})

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(b)
		case "POST":
			r = middleware.IsAuthenticated(w, r)

			if r == nil {
				return
			}

			userID := int(r.Context().Value("user").(jwt.MapClaims)["id"].(float64))

			decoder := json.NewDecoder(r.Body)
			var reqBody types.PostArticlesBody
			err := decoder.Decode(&reqBody)
			if err != nil {
				utils.SetErrorResponse(w, 400, err.Error())
				return
			}
			defer r.Body.Close()

			if reqBody.Title == "" {
				utils.SetErrorResponse(w, 400, "Title is required.")
				return
			}

			if reqBody.Body == "" {
				utils.SetErrorResponse(w, 400, "Body is required.")
				return
			}

			a := models.Article{
				UserID: userID,
				Title:  reqBody.Title,
				Body:   reqBody.Body,
			}

			var id int
			err = models.CreateArticle(db, &a).Scan(&id)

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			m, err := json.Marshal(a)

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(m)
		default:
			utils.SetErrorResponse(w, 404, "Invalid route.")
		}
	})
}

func HandleSignin(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if r.Body == nil {
				utils.SetErrorResponse(w, 400, "Username and password are required.")
				return
			}

			decoder := json.NewDecoder(r.Body)
			var sb types.SigninBody
			err := decoder.Decode(&sb)

			if err != nil {
				utils.SetErrorResponse(w, 400, "Username and password are required.")
				return
			}

			defer r.Body.Close()

			var (
				id       int
				username string
				email    string
				password string
			)
			err = models.GetUserByUsername(db, sb.Username).Scan(&id, &username, &email, &password)

			if err != nil {
				utils.SetErrorResponse(w, 400, "Username is invalid.")
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(password), []byte(sb.Password))

			if err != nil {
				utils.SetErrorResponse(w, 400, "Password is invalid.")
				return
			}

			c := types.JWTClaims{
				id,
				jwt.StandardClaims{
					Issuer: "Composition",
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
			ss, err := token.SignedString(types.JWTSecret)

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			m, err := json.Marshal(map[string]string{"token": ss})

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(m)
		default:
			utils.SetErrorResponse(w, 404, "Invalid route.")
		}
	}
}

func HandleSignup(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if r.Body == nil {
				utils.SetErrorResponse(w, 400, "Username, email, password, and password confirm are required.")
				return
			}

			decoder := json.NewDecoder(r.Body)
			var body types.SignupBody
			err := decoder.Decode(&body)
			if err != nil {
				utils.SetErrorResponse(w, 400, "Username, email, password, and password confirm are required.")
				return
			}

			defer r.Body.Close()

			if body == (types.SignupBody{}) {
				utils.SetErrorResponse(w, 400, "Username, email, password, and password confirm are required.")
				return
			}

			if body.Username == "" {
				utils.SetErrorResponse(w, 400, "Username is required.")
				return
			}

			if body.Email == "" {
				utils.SetErrorResponse(w, 400, "Email is required.")
				return
			}

			if body.Password == "" {
				utils.SetErrorResponse(w, 400, "Password is required.")
				return
			}

			if body.PasswordConfirm == "" {
				utils.SetErrorResponse(w, 400, "Password confirm is required.")
				return
			}

			if body.Password != body.PasswordConfirm {
				utils.SetErrorResponse(w, 400, "Passwords do not match.")
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			u := models.User{
				Username: body.Username,
				Email:    body.Email,
				Password: string(hash),
			}

			var id int
			err = models.CreateUser(db, &u).Scan(&id)

			if err != nil {
				utils.SetErrorResponse(w, 409, err.Error())
				return
			}

			c := types.JWTClaims{
				id,
				jwt.StandardClaims{
					Issuer: "Composition",
				},
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
			ss, err := token.SignedString(types.JWTSecret)

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			m, err := json.Marshal(map[string]string{"token": ss})

			if err != nil {
				utils.SetErrorResponse(w, 500, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(m)
		default:
			utils.SetErrorResponse(w, 404, "Invalid route.")
		}
	}
}
