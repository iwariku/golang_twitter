package controller

import (
	"golang_twitter/db"
	"log"
	"net/http"
	"strconv"

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

	// 1. Cookieを取得
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		log.Printf("Cookieの取得に失敗しました: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookieの取得に失敗しました"})
		return
	}

	// 2. RedisからUserIDを取得
	userIDStr, err := tc.Redis.Get(c.Request.Context(), sessionID).Result()
	if err != nil {
		log.Printf("セッション切れです。ログインしてください")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "セッション切れです。ログインしてください"})
		return
	}

	// 3. 型変換(DB保存の型と合わせるため)
	tempID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("string型からint型への変換が失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部でエラーが発生しました"})
		return
	}
	userID := int32(tempID)

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
