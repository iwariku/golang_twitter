package main

import (
	"context"
	"golang_twitter/controller"
	"golang_twitter/db"
	"golang_twitter/infrastructure"
	"golang_twitter/middleware"
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
	tc := &controller.TweetController{Queries: queries, Redis: redisClient}
	am := &middleware.AuthMiddleware{Redis: redisClient}

	r := gin.Default()
	r.GET("/health_check", func(c *gin.Context) {
		// JSONレスポンスを返す
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})

	r.LoadHTMLGlob("view/*")
	r.Static("/static", "./static")

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

	// グループを作成し、ミドルウェアを登録。
	authGroup := r.Group("/")
	authGroup.Use(am.CheckLogin)
	{
		authGroup.POST("/post", tc.TweetPost)

		authGroup.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", nil)
		})
		authGroup.GET("/api/tweets", tc.GetTweets)

		authGroup.GET("/tweet-detail", func(c *gin.Context) {
			c.HTML(http.StatusOK, "post-detail.html", nil)
		})
		authGroup.GET("/api/tweet-detail", tc.GetTweet)

		authGroup.GET("/user-detail", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user-detail.html", nil)
		})
		authGroup.GET("/api/user-detail", uc.GetUser)
		authGroup.GET("/api/user-tweets", uc.GetTweetsByUserID)
	}

	r.Run()
}
