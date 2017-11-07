package seeder

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/richardpanda/composition/server/api/models"
	"golang.org/x/crypto/bcrypt"
)

type article struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type user struct {
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Articles []article `json:"articles"`
}

func PopulateDB(db *sql.DB) {
	f, err := os.Open("./seeder/users.json")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}

	var users []user
	json.Unmarshal(b, &users)

	for _, user := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)

		if err != nil {
			log.Fatal(err)
		}

		u := &models.User{
			Username: user.Username,
			Email:    user.Email,
			Password: string(hash),
		}

		var id int
		models.CreateUser(db, u).Scan(&id)

		for _, article := range user.Articles {
			a := &models.Article{
				UserID: id,
				Title:  article.Title,
				Body:   article.Body,
			}

			models.CreateArticle(db, a)
		}
	}
}
