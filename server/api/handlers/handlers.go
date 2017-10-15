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
			if r.Body == nil {
				setResponse(w, 400, "Username and password are required.")
				return
			}

			decoder := json.NewDecoder(r.Body)
			var sb types.SigninBody
			err := decoder.Decode(&sb)

			if err != nil {
				setResponse(w, 400, "Username and password are required.")
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
				setResponse(w, 400, "Username is invalid.")
				return
			}

			err = bcrypt.CompareHashAndPassword([]byte(password), []byte(sb.Password))

			if err != nil {
				setResponse(w, 400, "Password is invalid.")
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
				setResponse(w, 500, err.Error())
				return
			}

			m, err := json.Marshal(map[string]string{"token": ss})

			if err != nil {
				setResponse(w, 500, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(m)
		default:
			setResponse(w, 404, "Invalid route.")
		}
	}
}

func HandleSignup(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			if r.Body == nil {
				setResponse(w, 400, "Username, email, password, and password confirm are required.")
				return
			}

			decoder := json.NewDecoder(r.Body)
			var body types.SignupBody
			err := decoder.Decode(&body)
			if err != nil {
				setResponse(w, 400, "Username, email, password, and password confirm are required.")
				return
			}

			defer r.Body.Close()

			if body == (types.SignupBody{}) {
				setResponse(w, 400, "Username, email, password, and password confirm are required.")
				return
			}

			if body.Username == "" {
				setResponse(w, 400, "Username is required.")
				return
			}

			if body.Email == "" {
				setResponse(w, 400, "Email is required.")
				return
			}

			if body.Password == "" {
				setResponse(w, 400, "Password is required.")
				return
			}

			if body.PasswordConfirm == "" {
				setResponse(w, 400, "Password confirm is required.")
				return
			}

			if body.Password != body.PasswordConfirm {
				setResponse(w, 400, "Passwords do not match.")
				return
			}

			hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)

			if err != nil {
				setResponse(w, 500, err.Error())
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
				setResponse(w, 409, err.Error())
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
				setResponse(w, 500, err.Error())
				return
			}

			m, err := json.Marshal(map[string]string{"token": ss})

			if err != nil {
				setResponse(w, 500, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write(m)
		default:
			setResponse(w, 404, "Invalid route.")
		}
	}
}

func setResponse(w http.ResponseWriter, code int, message string) {
	respBody, _ := json.Marshal(types.ErrorResponseBody{Message: message})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBody)
}
