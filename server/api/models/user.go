package models

import (
	"database/sql"
)

type User struct {
	Username string
	Email    string
	Password string
}

const createUserQuery = "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id;"
const createUsersTableQuery = `
	CREATE TABLE IF NOT EXISTS users (
		id       SERIAL       PRIMARY KEY,
		username VARCHAR(20)  UNIQUE NOT NULL,
		email    VARCHAR(50)  UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL
	);
`
const dropUserTableQuery = "DROP TABLE users CASCADE;"
const getUserByUsernameQuery = "SELECT * FROM users WHERE username=$1;"

func CreateUser(db *sql.DB, u *User) *sql.Row {
	return db.QueryRow(createUserQuery, u.Username, u.Email, u.Password)
}

func CreateUsersTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(createUsersTableQuery)
}

func DropUsersTable(db *sql.DB) (sql.Result, error) {
	return db.Exec(dropUserTableQuery)
}

func GetUserByUsername(db *sql.DB, username string) *sql.Row {
	return db.QueryRow(getUserByUsernameQuery, username)
}
