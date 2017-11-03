package controllers

import (
	"database/sql"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
)

func GetArticles(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	rows, err := models.GetLatestArticlePreviews(db, page)

	if err != nil {
		c.JSON(500, gin.H{"message": err.Error()})
		return
	}

	defer rows.Close()

	articlePreviews := []types.ArticlePreview{}

	for rows.Next() {
		var (
			username, title string
			id              int
			createdAt       time.Time
		)

		if err := rows.Scan(&username, &title, &id, &createdAt); err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		articlePreviews = append(articlePreviews, types.ArticlePreview{username, title, id, createdAt})
	}

	c.JSON(200, gin.H{"article_previews": articlePreviews})
}

func PostArticles(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)
	user, _ := c.Get("user")
	userID := int(user.(jwt.MapClaims)["id"].(float64))

	body := &types.PostArticlesBody{}

	if err := c.BindJSON(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if body.Title == "" {
		c.JSON(400, gin.H{"message": "Title is required."})
		return
	}

	if body.Body == "" {
		c.JSON(400, gin.H{"message": "Body is required."})
		return
	}

	a := &models.Article{
		UserID: userID,
		Title:  body.Title,
		Body:   body.Body,
	}

	var id int
	_ = models.CreateArticle(db, a).Scan(&id)

	c.JSON(201, gin.H{"article_id": id, "title": body.Title, "body": body.Body})
}
