package middleware

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// ここを認可をするためのリクエストの構造体にする
type AuthMiddleware struct {
	Redis *redis.Client
}

func (am *AuthMiddleware) CheckLogin(c *gin.Context) {
	// 1. Cookieを取得
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		log.Printf("Cookieの取得に失敗しました: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Cookieの取得に失敗しました"})
		c.Abort()
		return
	}

	// 2. RedisからUserIDを取得
	userIDStr, err := am.Redis.Get(c.Request.Context(), sessionID).Result()
	if err != nil {
		log.Printf("セッション切れです。ログインしてください")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "セッション切れです。ログインしてください"})
		c.Abort()
		return
	}

	// 3. 型変換(DB保存の型と合わせるため)
	tempID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("string型からint型への変換が失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部でエラーが発生しました"})
		c.Abort()
		return
	}
	userID := int32(tempID)

	c.Set("current_user_id", userID)
	c.Next()
}
