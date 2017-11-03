package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/richardpanda/composition/server/api/controllers"
	"github.com/richardpanda/composition/server/api/middleware"
)

func New(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.DB(db))

	r.GET("/api/articles", controllers.GetArticles)
	r.POST("/api/signin", controllers.PostSignin)
	r.POST("/api/signup", controllers.PostSignup)

	r.Use(middleware.Authenticate())

	r.POST("/api/articles", controllers.PostArticles)

	return r
}
