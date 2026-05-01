package main

import (
	"context"
	"github.com/iwariku/golang_twitter/controller"
	"github.com/iwariku/golang_twitter/db"
	"github.com/iwariku/golang_twitter/infrastructure"
	"github.com/iwariku/golang_twitter/middleware"
	"github.com/iwariku/golang_twitter/services/auth"
	"log"
	"net/http"
	"os"

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

	r.POST("/api/signup", uc.SignUp)
	r.GET("/api/activate", uc.Activate)
	r.POST("/api/login", uc.Login)

	// グループを作成し、ミドルウェアを登録。
	authGroup := r.Group("/")
	authGroup.Use(am.CheckLogin)
	{

		authGroup.GET("/api/users/me", uc.GetLoggedUserID)

		authGroup.POST("/api/post", tc.TweetPost)

		authGroup.GET("/api/tweets", tc.GetTweets)

		authGroup.GET("/api/tweets/:id", tc.GetTweet)

		authGroup.GET("/api/users/:id", uc.GetUser)
		authGroup.GET("/api/users/:id/tweets", uc.GetTweetsByUserID)

		authGroup.GET("/api/users/:id/retweets", tc.GetRetweetedTweetsByUserID)

		authGroup.GET("/api/user/bookmarks", tc.GetBookmarkedTweetsByUserID)

		authGroup.POST("/api/tweets/:id/like", tc.CreateLike)
		authGroup.DELETE("/api/tweets/:id/like", tc.DeleteLike)

		authGroup.POST("/api/tweets/:id/retweet", tc.CreateRetweet)
		authGroup.DELETE("/api/tweets/:id/retweet", tc.DeleteRetweet)

		authGroup.POST("/api/tweets/:id/bookmark", tc.CreateBookmark)
		authGroup.DELETE("/api/tweets/:id/bookmark", tc.DeleteBookmark)

		authGroup.POST("/api/users/:id/follow", uc.CreateFollow)
		authGroup.DELETE("/api/users/:id/follow", uc.DeleteFollow)

		authGroup.GET("/api/users/:id/followings", uc.GetFollowings)

		authGroup.GET("/api/users/:id/followers", uc.GetFollowers)

		authGroup.POST("/api/dm/group", dc.CreateGroup)

		authGroup.GET("/api/dm/groups", dc.GetGroups)

		authGroup.POST("/api/dm/add-member", dc.AddMemberToGroup)

		authGroup.POST("/api/dm/groups/:id/message", dc.CreateMessage)

		authGroup.GET("/api/dm/groups/:id/messages", dc.GetMessagesByGroupID)

		authGroup.DELETE("/api/user/unsubscribe", uc.DeleteUser)

	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
