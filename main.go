package main

import (
	"context"
	"golang_twitter/controller"
	"golang_twitter/db"
	"golang_twitter/infrastructure"
	"golang_twitter/middleware"
	"golang_twitter/services/auth"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".envが読み込めません")
	}

	ctx := context.Background()

	conn, queries := db.ConnectDB(ctx)
	defer conn.Close(ctx)

	// docker-compose.yamlの設定値をprodまたは、devと合わせることで使用できる(片方はコメントアウト済み)
	mailer := auth.NewProdMailer()
	// mailer := auth.NewDevMailer()

	redisClient := infrastructure.NewRedisClient()
	uc := &controller.UserController{Queries: queries, Mailer: mailer, Redis: redisClient}
	tc := &controller.TweetController{Queries: queries, Redis: redisClient}
	dc := &controller.DmController{Queries: queries, Redis: redisClient}
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

		authGroup.GET("/api/users/me", uc.GetLoggedUserID)

		authGroup.POST("/post", tc.TweetPost)

		authGroup.GET("/home", func(c *gin.Context) {
			c.HTML(http.StatusOK, "home.html", nil)
		})
		authGroup.GET("/api/tweets", tc.GetTweets)

		authGroup.GET("/tweet-detail/:id", func(c *gin.Context) {
			c.HTML(http.StatusOK, "post-detail.html", nil)
		})
		authGroup.GET("/api/tweets/:id", tc.GetTweet)

		authGroup.GET("/user-detail/:id", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user-detail.html", nil)
		})
		authGroup.GET("/api/users/:id", uc.GetUser)
		authGroup.GET("/api/users/:id/tweets", uc.GetTweetsByUserID)
		authGroup.GET("/user-retweet/:id", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user-retweet.html", nil)
		})
		authGroup.GET("/api/users/:id/retweets", tc.GetRetweetedTweetsByUserID)

		authGroup.GET("/user-bookmarks", func(c *gin.Context) {
			c.HTML(http.StatusOK, "user-bookmarks.html", nil)
		})
		authGroup.GET("/api/user/bookmarks", tc.GetBookmarkedTweetsByUserID)

		authGroup.POST("/api/tweets/:id/like", tc.CreateLike)
		authGroup.DELETE("/api/tweets/:id/like", tc.DeleteLike)

		authGroup.POST("/api/tweets/:id/retweet", tc.CreateRetweet)
		authGroup.DELETE("/api/tweets/:id/retweet", tc.DeleteRetweet)

		authGroup.POST("/api/tweets/:id/bookmark", tc.CreateBookmark)
		authGroup.DELETE("/api/tweets/:id/bookmark", tc.DeleteBookmark)

		authGroup.POST("/api/users/:id/follow", uc.CreateFollow)
		authGroup.DELETE("/api/users/:id/follow", uc.DeleteFollow)

		authGroup.GET("/users/:id/followings", func(c *gin.Context) {
			c.HTML(http.StatusOK, "follows.html", nil)
		})
		authGroup.GET("/api/users/:id/followings", uc.GetFollowings)

		authGroup.GET("/users/:id/followers", func(c *gin.Context) {
			c.HTML(http.StatusOK, "follows.html", nil)
		})
		authGroup.GET("/api/users/:id/followers", uc.GetFollowers)

		// --- DM機能 ---
		authGroup.GET("/dm/group", func(c *gin.Context) {
			c.HTML(http.StatusOK, "create-group.html", nil)
		})
		authGroup.POST("/api/dm/group", dc.CreateGroup)

		authGroup.GET("/dm/groups", func(c *gin.Context) {
			c.HTML(http.StatusOK, "groups.html", nil)
		})
		authGroup.GET("/api/dm/groups", dc.GetGroups)

		authGroup.GET("/dm/add-member", func(c *gin.Context) {
			c.HTML(http.StatusOK, "add-member.html", nil)
		})
		authGroup.POST("/api/dm/add-member", dc.AddMemberToGroup)

		authGroup.GET("dm/groups/:id/message", func(c *gin.Context) {
			c.HTML(http.StatusOK, "message.html", nil)
		})
		authGroup.POST("/api/dm/groups/:id/message", dc.CreateMessage)

		authGroup.GET("/dm/groups/:id/messages", func(c *gin.Context) {
			c.HTML(http.StatusOK, "groups-messages.html", nil)
		})
		authGroup.GET("/api/dm/groups/:id/messages", dc.GetMessagesByGroupID)

		authGroup.GET("/user/unsubscribe", func(c *gin.Context) {
			c.HTML(http.StatusOK, "unsubscribe.html", nil)
		})
		authGroup.DELETE("/api/user/unsubscribe", uc.DeleteUser)

	}

	r.Run()
}
