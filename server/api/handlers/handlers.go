package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
	"golang.org/x/crypto/bcrypt"
)

func HandlePostArticles(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			userID := int(r.Context().Value("user").(jwt.MapClaims)["id"].(float64))

			decoder := json.NewDecoder(r.Body)
			var reqBody types.PostArticlesBody
			err := decoder.Decode(&reqBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			defer r.Body.Close()

			a := models.Article{
				UserID: userID,
				Title:  reqBody.Title,
				Body:   reqBody.Body,
			}

			var id int
			err = models.CreateArticle(db, &a).Scan(&id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

			m, err := json.Marshal(a)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			w.Write(m)
		default:
			http.Error(w, "Invalid route.", http.StatusNotFound)
		}
	})
}

func HandleSignin(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			var sb types.SigninBody
			err := decoder.Decode(&sb)

			if err != nil {
				http.Error(w, "Username and password are required.", http.StatusBadRequest)
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
				http.Error(w, "Username is invalid.", http.StatusBadRequest)
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(password), []byte(sb.Password))

			if err != nil {
				http.Error(w, "Password is invalid.", http.StatusBadRequest)
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			m, err := json.Marshal(map[string]string{"token": ss})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(m)
		default:
			http.Error(w, "Invalid route.", http.StatusNotFound)
		}
	}
}

func HandleSignup(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			decoder := json.NewDecoder(r.Body)
			var sb types.SignupBody

			err := decoder.Decode(&sb)

			if err != nil {
				http.Error(w, "Username, email, password, and password confirm are required.", http.StatusBadRequest)
				panic(err)
			}

			defer r.Body.Close()

			if sb.Password != sb.PasswordConfirm {
				http.Error(w, "Passwords do not match.", http.StatusBadRequest)
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(sb.Password), bcrypt.MinCost)

			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			u := models.User{
				Username: sb.Username,
				Email:    sb.Email,
				Password: string(hash),
			}

			var id int
			err = models.CreateUser(db, &u).Scan(&id)

			if err != nil {
				http.Error(w, err.Error(), http.StatusConflict)
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			m, err := json.Marshal(map[string]string{"token": ss})

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(m)
		default:
			http.Error(w, "Invalid route.", http.StatusNotFound)
		}
	}
}
