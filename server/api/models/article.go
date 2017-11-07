package models

import (
	"database/sql"
)

type Article struct {
	UserID int    `json:"id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

const createArticleQuery = "INSERT INTO articles (user_id, title, body, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id;"
const createArticlesTableQuery = `
	CREATE TABLE IF NOT EXISTS articles (
		id         SERIAL       PRIMARY KEY,
		user_id    SERIAL       REFERENCES users,
		title      VARCHAR(100)  NOT NULL,
		body       TEXT         NOT NULL,
		created_at TIMESTAMP    NOT NULL
	);
`
const dropArticlesTableQuery = "DROP TABLE articles;"
const getLatestArticlePreviewsQuery = `
	SELECT username, title, articles.id, created_at
	FROM users, articles
	WHERE users.id = articles.user_id
	ORDER BY created_at DESC
	LIMIT 10
	OFFSET $1;
`

func CreateArticle(db *sql.DB, a *Article) *sql.Row {
	return db.QueryRow(createArticleQuery, a.UserID, a.Title, a.Body)
}

func CreateArticlesTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(createArticlesTableQuery)
}

func DropArticlesTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(dropArticlesTableQuery)
}

func GetLatestArticlePreviews(db *sql.DB, page int) (*sql.Rows, error) {
	return db.Query(getLatestArticlePreviewsQuery, (page-1)*10)
}
