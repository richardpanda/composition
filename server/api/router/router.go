package router

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/richardpanda/composition/server/api/controllers"
	"github.com/richardpanda/composition/server/api/middlewares"
)

func New(db *sql.DB) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.DB(db))

	r.GET("/api/articles/:id", controllers.GetArticle)
	r.GET("/api/articles", controllers.GetArticles)
	r.POST("/api/signin", controllers.PostSignin)
	r.POST("/api/signup", controllers.PostSignup)

	r.Use(middlewares.Authenticate())

	r.POST("/api/articles", controllers.PostArticles)

	return r
}
