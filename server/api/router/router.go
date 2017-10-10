package router

import (
	"database/sql"
	"net/http"

	"github.com/richardpanda/composition/server/api/handlers"
	"github.com/richardpanda/composition/server/api/middleware"
)

func New(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/api/articles", middleware.IsAuthenticated(handlers.HandlePostArticles(db)))
	mux.HandleFunc("/api/signin", handlers.HandleSignin(db))
	mux.HandleFunc("/api/signup", handlers.HandleSignup(db))
	return mux
}
