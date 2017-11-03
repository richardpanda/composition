package controllers

import (
	"database/sql"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/richardpanda/composition/server/api/models"
	"github.com/richardpanda/composition/server/api/types"
	"golang.org/x/crypto/bcrypt"
)

func PostSignin(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	if c.Request.Body == nil {
		c.JSON(400, gin.H{"message": "Username and password are required."})
		return
	}

	body := &types.SigninBody{}

	if err := c.BindJSON(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	var (
		id       int
		username string
		email    string
		password string
	)

	err := models.GetUserByUsername(db, body.Username).Scan(&id, &username, &email, &password)

	if err != nil {
		c.JSON(400, gin.H{"message": "Username is invalid."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(body.Password))

	if err != nil {
		c.JSON(400, gin.H{"message": "Password is invalid."})
		return
	}

	claims := types.JWTClaims{
		id,
		jwt.StandardClaims{
			Issuer: "Composition",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(types.JWTSecret)

	c.JSON(200, gin.H{"token": ss})
}

func PostSignup(c *gin.Context) {
	db := c.MustGet("db").(*sql.DB)

	if c.Request.Body == nil {
		c.JSON(400, gin.H{"message": "Username, email, password, and password confirm are required."})
		return
	}

	body := &types.SignupBody{}

	if err := c.BindJSON(body); err != nil {
		c.JSON(400, gin.H{"message": err.Error()})
		return
	}

	if body.Username == "" {
		c.JSON(400, gin.H{"message": "Username is required."})
		return
	}

	if body.Email == "" {
		c.JSON(400, gin.H{"message": "Email is required."})
		return
	}

	if body.Password == "" {
		c.JSON(400, gin.H{"message": "Password is required."})
		return
	}

	if body.PasswordConfirm == "" {
		c.JSON(400, gin.H{"message": "Password confirm is required."})
		return
	}

	if body.Password != body.PasswordConfirm {
		c.JSON(400, gin.H{"message": "Passwords do not match."})
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)

	u := &models.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(hash),
	}

	var id int
	err := models.CreateUser(db, u).Scan(&id)

	if err, ok := err.(*pq.Error); ok {
		fmt.Printf("Severity: %s\n", err.Severity)
		fmt.Printf("Code: %s\n", err.Code)
		fmt.Printf("Message: %s\n", err.Message)
		fmt.Printf("Detail: %s\n", err.Detail)
		fmt.Printf("Hint: %s\n", err.Hint)
		fmt.Printf("Position: %s\n", err.Position)
		fmt.Printf("InternalPosition: %s\n", err.InternalPosition)
		fmt.Printf("InternalQuery: %s\n", err.InternalQuery)
		fmt.Printf("Where: %s\n", err.Where)
		fmt.Printf("Schema: %s\n", err.Schema)
		fmt.Printf("Table: %s\n", err.Table)
		fmt.Printf("Column: %s\n", err.Column)
		fmt.Printf("DataTypeName: %s\n", err.DataTypeName)
		fmt.Printf("Constraint: %s\n", err.Constraint)
		fmt.Printf("File: %s\n", err.File)
		fmt.Printf("Line: %s\n", err.Line)
		fmt.Printf("Routine: %s\n", err.Routine)

		switch err.Constraint {
		case "users_username_key":
			c.JSON(400, gin.H{"message": "Username is not available."})
		case "users_email_key":
			c.JSON(400, gin.H{"message": "Email is not available."})
		default:
			fmt.Println(err)
			c.JSON(500, gin.H{"message": err.Error()})
		}
		return
	}

	claims := types.JWTClaims{
		id,
		jwt.StandardClaims{
			Issuer: "Composition",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(types.JWTSecret)

	c.JSON(200, gin.H{"token": ss})
}
