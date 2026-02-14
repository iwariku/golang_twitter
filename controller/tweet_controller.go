package controller

import (
	"golang_twitter/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type TweetController struct {
	Queries *db.Queries
	Redis   *redis.Client
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

	// リクエストスコープに保存されたcurrent_user_idを取得
	userIDAny, exists := c.Get("current_user_id")
	if !exists {
		log.Printf("ログインが必要です")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	// 型変換チェック(anyをint32として証明するため)
	userID, ok := userIDAny.(int32)
	if !ok {
		log.Printf("リクエストスコープ内のUserIDがint32ではありませんでした")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "内部エラーが発生しました"})
		return
	}

	tweet, err := tc.Queries.CreateTweet(c.Request.Context(), db.CreateTweetParams{
		UserID:  userID,
		Content: req.Content,
	})
	if err != nil {
		log.Printf("DBへの保存に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBへの保存に失敗しました"})
		return
	}
	c.JSON(http.StatusCreated, tweet)
}
