package router

import (
	"database/sql"
	"net/http"

	"github.com/richardpanda/composition/server/api/handlers"
)

func New(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/signup", handlers.HandleSignup(db))
	return mux
}