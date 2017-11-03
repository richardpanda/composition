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

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	_ = models.CreateUser(db, u)

	b, _ := json.Marshal(types.SigninBody{
		Username: "test",
		Password: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signin", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 200)
	assertJSONHeader(t, rr)

	respBody := &types.SigninResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
}

func TestSigninWithoutBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/signin", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Username and password are required.")
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
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Username is invalid.")
}

func TestSigninWithInvalidPassword(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	models.CreateUser(db, u)

	b, _ := json.Marshal(types.SigninBody{
		Username: "test",
		Password: "invalid password",
	})

	req, _ := http.NewRequest("POST", "/api/signin", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Password is invalid.")
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
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 200)
	assertJSONHeader(t, rr)

	respBody := &types.SigninResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
}

func TestSignupWithoutBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/signup", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Username, email, password, and password confirm are required.")
}

func TestSignupWithEmptyBody(t *testing.T) {
	b, _ := json.Marshal(types.SignupBody{})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Username is required.")
}

func TestSignupWithoutUsername(t *testing.T) {
	b, _ := json.Marshal(types.SignupBody{
		Email:           "test@test.com",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Username is required.")
}

func TestSignupWithoutEmail(t *testing.T) {
	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Email is required.")
}

func TestSignupWithoutPassword(t *testing.T) {
	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Email:           "test@test.com",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Password is required.")
}

func TestSignupWithoutPasswordConfirm(t *testing.T) {
	b, _ := json.Marshal(types.SignupBody{
		Username: "test",
		Email:    "test@test.com",
		Password: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Password confirm is required.")
}

func TestSignupWithMismatchPasswordAndPasswordConfirm(t *testing.T) {
	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Email:           "test@test.com",
		Password:        "a",
		PasswordConfirm: "b",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Passwords do not match.")
}

func TestSignUpWithRegisteredUsername(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	_ = models.CreateUser(db, u)

	b, _ := json.Marshal(types.SignupBody{
		Username:        "test",
		Email:           "test2@test.com",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Username is not available.")
}

func TestSignUpWithRegisteredEmail(t *testing.T) {
	createUsersTable()
	defer dropUsersTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	_ = models.CreateUser(db, u)

	b, _ := json.Marshal(types.SignupBody{
		Username:        "test2",
		Email:           "test@test.com",
		Password:        "test",
		PasswordConfirm: "test",
	})

	req, _ := http.NewRequest("POST", "/api/signup", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Email is not available.")
}
