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

func TestGetArticleWithExistentArticle(t *testing.T) {
	createUsersTable()
	createArticlesTable()
	defer dropUsersTable()
	defer dropArticlesTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	var userID int
	err = models.CreateUser(db, u).Scan(&userID)

	assertEqual(t, err, nil)

	a := &models.Article{
		UserID: userID,
		Title:  "Title",
		Body:   "Body",
	}

	var articleID int
	err = models.CreateArticle(db, a).Scan(&articleID)

	endpoint := fmt.Sprintf("/api/articles/%d", articleID)
	req, _ := http.NewRequest("GET", endpoint, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 200)
	assertJSONHeader(t, rr)

	r := &types.GetArticleResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), r)

	assertEqual(t, err, nil)
	assertEqual(t, r.ID, articleID)
	assertEqual(t, r.Title, "Title")
	assertEqual(t, r.Body, "Body")
	assertEqual(t, r.Username, "test")
}

func TestGetArticleWithNonexistentArticle(t *testing.T) {
	createUsersTable()
	createArticlesTable()
	defer dropUsersTable()
	defer dropArticlesTable()

	endpoint := fmt.Sprintf("/api/articles/1")
	req, _ := http.NewRequest("GET", endpoint, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 404)
	assertJSONHeader(t, rr)

	r := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), r)

	assertEqual(t, err, nil)
	assertEqual(t, r.Message, "Unable to find article.")
}

func TestGetArticlePreviews(t *testing.T) {
	createUsersTable()
	createArticlesTable()
	defer dropUsersTable()
	defer dropArticlesTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	var id int
	err = models.CreateUser(db, u).Scan(&id)

	assertEqual(t, err, nil)

	for i := 1; i <= 11; i++ {
		models.CreateArticle(db, &models.Article{
			UserID: id,
			Title:  fmt.Sprintf("Title %d", i),
			Body:   fmt.Sprintf("Body %d", i),
		})
	}

	req, _ := http.NewRequest("GET", "/api/articles", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 200)
	assertJSONHeader(t, rr)

	resp := &types.GetArticlesResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), resp)

	assertEqual(t, err, nil)
	assertEqual(t, len(resp.ArticlePreviews), 10)
	assertEqual(t, resp.ArticlePreviews[0].ID, 11)
	assertEqual(t, resp.ArticlePreviews[0].Title, "Title 11")
	assertEqual(t, resp.ArticlePreviews[0].Username, "test")
	assertEqual(t, resp.ArticlePreviews[9].ID, 2)
	assertEqual(t, resp.ArticlePreviews[9].Title, "Title 2")
	assertEqual(t, resp.ArticlePreviews[9].Username, "test")
}

func TestSuccessfulPostArticles(t *testing.T) {
	createUsersTable()
	createArticlesTable()
	defer dropUsersTable()
	defer dropArticlesTable()

	password := "test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	assertEqual(t, err, nil)

	u := &models.User{
		Username: "test",
		Email:    "test@test.com",
		Password: string(hash),
	}

	var id int
	err = models.CreateUser(db, u).Scan(&id)

	assertEqual(t, err, nil)

	c := types.JWTClaims{
		id,
		jwt.StandardClaims{
			Issuer: "Composition",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString(types.JWTSecret)

	assertEqual(t, err, nil)

	b, _ := json.Marshal(types.PostArticlesRequestBody{
		Title: "Lorem Ipsum",
		Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 201)
	assertJSONHeader(t, rr)

	respBody := &types.PostArticlesResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
}

func TestPostArticlesWithoutToken(t *testing.T) {
	b, _ := json.Marshal(types.PostArticlesRequestBody{
		Title: "Lorem Ipsum",
		Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 401)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err := json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Authorization header is required.")
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

	assertEqual(t, err, nil)

	b, _ := json.Marshal(types.PostArticlesRequestBody{
		Title: "Lorem Ipsum",
		Body:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Invalid token.")
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

	assertEqual(t, err, nil)

	b, _ := json.Marshal(types.PostArticlesRequestBody{
		Body: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Title is required.")
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

	assertEqual(t, err, nil)

	b, _ := json.Marshal(types.PostArticlesRequestBody{
		Title: "Lorem Ipsum",
	})

	authHeader := fmt.Sprintf("Bearer %s", ss)
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBuffer(b))
	req.Header.Set("Authorization", authHeader)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assertEqual(t, rr.Code, 400)
	assertJSONHeader(t, rr)

	respBody := &types.ErrorResponseBody{}
	err = json.Unmarshal(rr.Body.Bytes(), respBody)

	assertEqual(t, err, nil)
	assertEqual(t, respBody.Message, "Body is required.")
}
