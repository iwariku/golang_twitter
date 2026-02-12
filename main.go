package main

import (
	"context"
	"golang_twitter/controller"
	"golang_twitter/db"
	"golang_twitter/infrastructure"
	"golang_twitter/services/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	conn, queries := db.ConnectDB(ctx)
	defer conn.Close(ctx)

	mailer := auth.NewMailer()
	redisClient := infrastructure.NewRedisClient()
	uc := &controller.UserController{Queries: queries, Mailer: mailer, Redis: redisClient}
	tc := &controller.TweetController{Queries: queries}

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

	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
	r.POST("/login", uc.Login)

	r.GET("/activate", uc.Activate)

	r.GET("/post", func(c *gin.Context) {
		c.HTML(http.StatusOK, "post.html", nil)
	})
	r.POST("/post", tc.TweetPost)

	r.Run()
}
