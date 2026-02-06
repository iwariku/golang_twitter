package main

import (
	"context"
	"golang_twitter/controller"
	"golang_twitter/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	conn, queries := db.ConnectDB(ctx)
	defer conn.Close(ctx)

	uc := &controller.UserController{Queries: queries}

	r := gin.Default()
	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")

	r.GET("/health_check", func(c *gin.Context) {
		// JSONレスポンスを返す
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})

	r.POST("/signup", uc.SignUp)

	r.Run()
}
