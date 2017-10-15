package router

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/richardpanda/composition/server/api/models"
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
