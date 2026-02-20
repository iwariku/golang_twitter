package controller

import (
	"fmt"
	"golang_twitter/db"
	"golang_twitter/utils"
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

type TweetResponse struct {
	UserID  int32  `json:"user_id"`
	Content string `json:"content"`
}

type PaginatedTweetsResponse struct {
	Tweets []TweetResponse `json:"tweets"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Count  int             `json:"count"`
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

// utils/param.goでstring->int型の関数を作成したので、別ブランチでリファクタリングする。終わり次第コメントを消す
func (tc *TweetController) GetTweets(c *gin.Context) {
	// 1. URLパラメーターから文字列を取得
	limitStr := c.Query("limit")
	if limitStr == "" {
		limitStr = "10"
	}
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offsetStr = "0"
	}

	// 2. 文字列をint型に変換する
	limitInt, err := strconv.Atoi(limitStr)
	if err != nil {
		limitInt = 10
	}
	offsetInt, err := strconv.Atoi(offsetStr)
	if err != nil {
		offsetInt = 0
	}

	// 件数取得
	TotalCount, err := tc.Queries.GetTweetCount(c.Request.Context())
	if err != nil {
		log.Printf("件数取得に失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "件数取得に失敗失敗しました"})
		return
	}

	tweets, err := tc.Queries.GetTweets(c.Request.Context(), db.GetTweetsParams{
		Limit:  int32(limitInt),
		Offset: int32(offsetInt),
	})
	if err != nil {
		log.Printf("DBからの取得に失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBからの取得に失敗しました"})
		return
	}

	var TweetsRes []TweetResponse
	for _, t := range tweets {
		TweetsRes = append(TweetsRes, TweetResponse{
			UserID:  t.UserID,
			Content: t.Content,
		})
	}

	paginatedRes := PaginatedTweetsResponse{
		Tweets: TweetsRes,
		Limit:  limitInt,
		Offset: offsetInt,
		Count:  int(TotalCount),
	}

	c.JSON(http.StatusOK, paginatedRes)
}

func (tc *TweetController) GetTweet(c *gin.Context) {
	id, err := utils.ParseQueryInt32(c, "id")
	if err != nil {
		log.Printf("パラメータ解析に失敗しました: %v", err)
		c.JSON(http.StatusBadRequest, "不正なリクエストです")
		return
	}

	tweet, err := tc.Queries.GetTweet(c.Request.Context(), id)
	if err != nil {
		log.Printf("DBからの取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, "DBからの取得に失敗しました")
		return
	}

	TweetRes := TweetResponse{
		UserID:  tweet.UserID,
		Content: tweet.Content,
	}

	c.JSON(http.StatusOK, TweetRes)
}
