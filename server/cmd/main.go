package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/router"
	"github.com/richardpanda/composition/server/seeder"
)

func main() {
	env := os.Getenv("ENVIRONMENT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("user=%s dbname=%s sslmode=disable", user, dbname)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		log.Fatal(err)
	}

	_, err = models.CreateUsersTable(db)

	if err != nil {
		log.Fatal(err)
	}

	_, err = models.CreateArticlesTable(db)

	if err != nil {
		log.Fatal(err)
	}

	if env == "dev" {
		seeder.PopulateDB(db)
	}

	http.ListenAndServe(":8080", router.New(db))
}
