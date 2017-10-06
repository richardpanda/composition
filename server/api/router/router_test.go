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

	_ "github.com/lib/pq"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
)

var (
	user             = os.Getenv("DB_USER")
	dbname           = os.Getenv("TEST_DB_NAME")
	connectionString = fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname)
	db, _            = sql.Open("postgres", connectionString)
	mux              = New(db)
)

func setup() {
	_, err := models.CreateUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func teardown() {
	_, err := models.DropUsersTable(db)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSignup(t *testing.T) {
	setup()
	defer teardown()

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

	var sr types.SignupResponse
	err := json.Unmarshal(rr.Body.Bytes(), &sr)

	if err != nil {
		t.Fatal(err)
	}
}
