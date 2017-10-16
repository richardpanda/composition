package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
	"golang.org/x/crypto/bcrypt"
)

func TestSuccessfulPostArticles(t *testing.T) {
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

	expectedCode := 201
	if rr.Code != expectedCode {
		t.Fatalf("Expected: %d\nReceived: %d", expectedCode, rr.Code)
	}

	var resp models.Article
	err = json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		t.Fatal(err)
	}
}

func TestPostArticlesWithoutToken(t *testing.T) {
	b, _ := json.Marshal(types.PostArticlesBody{
		Title: "Lorem Ipsum",
		Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	expectedCode := 401
	if rr.Code != expectedCode {
		t.Fatalf("Expected: %d\nReceived: %d", expectedCode, rr.Code)
	}

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected: %s\nReceived %s", expectedContentType, rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Authorization header is required."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
	}
}

func TestPostArticlesWithInvalidToken(t *testing.T) {
	id := 1
	c := types.JWTClaims{
		id,
		jwt.StandardClaims{
			Issuer: "Composition",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	wrongSecret := append(types.JWTSecret, 'a')
	ss, err := token.SignedString(wrongSecret)
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

	expectedCode := 400
	if rr.Code != expectedCode {
		t.Fatalf("Expected: %d\nReceived: %d", expectedCode, rr.Code)
	}

	expectedContentType := "application/json"
	if rr.Header().Get("Content-Type") != expectedContentType {
		t.Fatalf("Expected: %s\nReceived %s", expectedContentType, rr.Header().Get("Content-Type"))
	}

	var respBody types.ErrorResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Invalid token."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
	}
}

func TestPostArticlesWithoutTitle(t *testing.T) {
	id := 1
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
		Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	expectedCode := 400
	if rr.Code != expectedCode {
		t.Fatalf("Expected: %d\nReceived: %d", expectedCode, rr.Code)
	}

	var resp models.Article
	err = json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		t.Fatal(err)
	}

	var respBody types.ErrorResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Title is required."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
	}
}

func TestPostArticlesWithoutBody(t *testing.T) {
	id := 1
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
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	expectedCode := 400
	if rr.Code != expectedCode {
		t.Fatalf("Expected: %d\nReceived: %d", expectedCode, rr.Code)
	}

	var resp models.Article
	err = json.Unmarshal(rr.Body.Bytes(), &resp)

	if err != nil {
		t.Fatal(err)
	}

	var respBody types.ErrorResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Body is required."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
	}
}
