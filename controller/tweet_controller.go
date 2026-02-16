package controller

import (
	"fmt"
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

type TweetResponse struct {
	UserID  int32  `json:"user_id"`
	Content string `json:"content"`
}

func GetUserIDFromContext(c *gin.Context) (int32, error) {
	// リクエストスコープに保存されたcurrent_user_idを取得
	userIDAny, exists := c.Get("current_user_id")
	if !exists {
		return 0, fmt.Errorf("user_idがコンテキストに設定されていません")
	}

	// 型変換チェック(anyをint32として証明するため)
	userID, ok := userIDAny.(int32)
	if !ok {
		return 0, fmt.Errorf("user_idはint32型ではありません")
	}
	return userID, nil
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

	userID, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
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

	TweetRes := TweetResponse{
		UserID:  tweet.UserID,
		Content: tweet.Content,
	}

	c.JSON(http.StatusCreated, TweetRes)
}
