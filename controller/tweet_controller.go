package controller

import (
	"golang_twitter/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	// "github.com/redis/go-redis/v9"
)

// PRがマージされたらredisのmodが入るけどこれは古いからまだ入っていない
// 2/11PRをマージできなかったらgo getしてもいいかも
type TweetController struct {
	Queries *db.Queries
	// Redis   *redis.Client
}

// Tweet投稿の流れ
// リクエストするユーザーを取得
// CreateTweetの引数に取得したユーザーとcontentを渡す
// フロントからのCookieがRedisにあるかを判定 trueだったら保存するロジックへ falseだったらreturn
func (tc *TweetController) TweetPost(c *gin.Context) {
	var req TweetRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSONの形式が違います: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSONの形式が違います"})
		return
	}

	tweet, err := tc.Queries.CreateTweet(c.Request.Context(), db.CreateTweetParams{
		UserID:  1,
		Content: req.Content,
	})
	if err != nil {
		log.Printf("DBへの保存に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBへの保存に失敗しました"})
		return
	}
	c.JSON(http.StatusCreated, tweet)
}
