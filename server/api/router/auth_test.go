package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
	"golang.org/x/crypto/bcrypt"
)

func TestSuccessfulSignin(t *testing.T) {
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

	expectedCode := 200
	if rr.Code != expectedCode {
		t.Fatalf("Expected: %d\nReceived: %d", expectedCode, rr.Code)
	}

	var respBody types.SigninResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}
}

func TestSigninWithoutBody(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	req, _ := http.NewRequest("POST", "/api/signin", nil)
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
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Username and password are required."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
	}
}

func TestSigninWithInvalidUsername(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	b, _ := json.Marshal(types.SigninBody{
		Username: "test",
		Password: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signin", bytes.NewBuffer(b))
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
	err := json.Unmarshal(rr.Body.Bytes(), &respBody)

	if err != nil {
		t.Fatal(err)
	}

	expectedMessage := "Username is invalid."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
	}
}

func TestSigninWithInvalidPassword(t *testing.T) {
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
		Password: "invalid password",
	})

	req, _ := http.NewRequest("POST", "/api/signin", bytes.NewBuffer(b))
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

	expectedMessage := "Password is invalid."
	if respBody.Message != expectedMessage {
		t.Fatalf("Expected: %s\nReceived: %s", expectedMessage, respBody.Message)
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
