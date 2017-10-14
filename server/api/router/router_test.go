package router

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	user             = os.Getenv("DB_USER")
	dbname           = os.Getenv("TEST_DB_NAME")
	connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname)
	db, _            = sql.Open("postgres", connectionString)
	mux              = New(db)
)

func createArticlesTable() {
	_, err := models.CreateArticlesTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func dropArticlesTable() {
	_, err := models.DropArticlesTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func createUsersTable() {
	_, err := models.CreateUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func dropUsersTable() {
	_, err := models.DropUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func TestPostArticles(t *testing.T) {
	createUsersTable()
	createArticlesTable()
	defer dropUsersTable()
	defer dropArticlesTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}

	u := models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	var id int
	err = models.CreateUser(db, &u).Scan(&id)
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	b, _ := json.Marshal(types.PostArticlesBody{
		Title: "Lorem Ipsum",
		Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 201 {
		fmt.Println(rr)
		t.Fatalf("Expected: 201, received: %d", rr.Code)
	}

	var resp models.Article
	err = json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSignin(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		t.Fatal(err)
	}

	u := models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	models.CreateUser(db, &u)

	b, _ := json.Marshal(types.SigninBody{
		Username: "test",
		Password: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signin", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		fmt.Println(rr)
		t.Fatalf("Expected: 200, received: %d", rr.Code)
	}

	var sr types.SigninResponse
	err = json.Unmarshal(rr.Body.Bytes(), &sr)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSuccessfulSignup(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Email:           "test@test.com",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("Expected: 200, received: %d", rr.Code)
	}

	var respBody types.SignupResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSignupWithoutBody(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	req, _ := http.NewRequest("POST", "/api/signup", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Username, email, password, and password confirm are required." {
		t.Fatalf("Expected: Username, email, password, and password confirm are required.\nReceived: %s", respBody.Message)
	}
}

func TestSignupWithEmptyBody(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Username, email, password, and password confirm are required." {
		t.Fatalf("Expected: Username, email, password, and password confirm are required.\nReceived: %s", respBody.Message)
	}
}

func TestSignupWithoutUsername(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{
		Email:           "test@test.com",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Username is required." {
		t.Fatalf("Expected: Username is required.\nReceived: %s", respBody.Message)
	}
}

func TestSignupWithoutEmail(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Email is required." {
		t.Fatalf("Expected: Email is required.\nReceived: %s", respBody.Message)
	}
}

func TestSignupWithoutPassword(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Email:           "test@test.com",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Password is required." {
		t.Fatalf("Expected: Password is required.\nReceived: %s", respBody.Message)
	}
}

func TestSignupWithoutPasswordConfirm(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Password confirm is required." {
		t.Fatalf("Expected: Password confirm is required.\nReceived: %s", respBody.Message)
	}
}

func TestSignupWithMismatchPasswordAndPasswordConfirm(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Email:           "test@test.com",
		Password:        "a",
		PasswordConfirm: "b",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != 400 {
		t.Fatalf("Expected: 400\nReceived: %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json" {
		t.Fatalf("Expected: application/json\nReceived %s", rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	if respBody.Message != "Passwords do not match." {
		t.Fatalf("Expected: Passwords do not match.\nReceived: %s", respBody.Message)
	}
}
