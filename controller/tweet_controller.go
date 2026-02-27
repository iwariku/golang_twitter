package controller

import (
	"fmt"
	"golang_twitter/db"
	"golang_twitter/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type TweetController struct {
	Queries *db.Queries
	Redis   *redis.Client
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

	tweetRes := FormatTweetResponse(tweet)

	c.JSON(http.StatusCreated, tweetRes)
}

// func (tc *TweetController) GetTweets(c *gin.Context) {
// 	limit, err := utils.ParseQueryInt32WithDefault(c, "limit", 10)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "limitの形式が違います"})
// 		return
// 	}

// 	offset, err := utils.ParseQueryInt32WithDefault(c, "offset", 0)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "offsetの形式が違います"})
// 		return
// 	}

// 	// 件数取得
// 	totalCount, err := tc.Queries.GetTweetCount(c.Request.Context())
// 	if err != nil {
// 		log.Printf("件数取得に失敗しました")
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "件数取得に失敗失敗しました"})
// 		return
// 	}

// 	tweets, err := tc.Queries.GetTweets(c.Request.Context(), db.GetTweetsParams{
// 		Limit:  limit,
// 		Offset: offset,
// 	})
// 	if err != nil {
// 		log.Printf("DBからの取得に失敗: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBからの取得に失敗しました"})
// 		return
// 	}

// 	paginatedTweetsResponse := FormatPaginatedTweetsResponse(tweets, limit, offset, totalCount)

// 	c.JSON(http.StatusOK, paginatedTweetsResponse)
// }

// ===================
// いいねツイート一覧機能
// ===================
// 最終的にこちらを使用する。(問題がなければGetTweetsに命名変更)
func (tc *TweetController) GetTweetsWithLikes(c *gin.Context) {
	//ログインユーザーの取得はここだろうと予想
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	limit, err := utils.ParseQueryInt32WithDefault(c, "limit", 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limitの形式が違います"})
		return
	}

	offset, err := utils.ParseQueryInt32WithDefault(c, "offset", 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offsetの形式が違います"})
		return
	}

	totalCount, err := tc.Queries.GetTweetCount(c.Request.Context())
	if err != nil {
		log.Printf("件数取得に失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "件数取得に失敗失敗しました"})
		return
	}

	dbTweets, err := tc.Queries.GetTweetsWithLikes(c.Request.Context(), db.GetTweetsWithLikesParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBからの取得に失敗しました"})
		return
	}

	paginatedTweetsRes := FormatPaginatedWithLikeTweetsResponse(dbTweets, limit, offset, totalCount)

	c.JSON(http.StatusOK, paginatedTweetsRes)
}

func (tc *TweetController) GetTweet(c *gin.Context) {
	id, err := utils.ParseParamInt32(c, "id")
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

	tweetRes := FormatTweetResponse(tweet)

	c.JSON(http.StatusOK, tweetRes)
}

// いいね機能
// バックエンドでツイートの有無の状態を確認する

// 1. 今何のツイートとユーザーIDなのかを確認する
// 2. 1.の情報を使いDBに該当するツイートがあるかを確認する
// 3. 2.の結果を元に条件分岐でcreateLikeかdeleteLikeかを決める
// 4. データを整形してレスポンスを返す
func (tc *TweetController) ToggleLike(c *gin.Context) {
	var currentLikeStatus bool

	// userId, err := GetUserIDFromContext(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
	// 	return
	// }

	// ===========================
	// 一旦userIdを1と固定して進める
	// ===========================
	userId := int32(1)

	tweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	// Likeを持つ == DBにレコードがある
	hasLiked, err := tc.Queries.GetLikeExists(c.Request.Context(), db.GetLikeExistsParams{
		UserID:  userId,
		TweetID: tweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "条件に合致するツイートがDBにありません"})
		return
	}

	if hasLiked {
		err := tc.Queries.DeleteLike(c.Request.Context(), db.DeleteLikeParams{
			UserID:  userId,
			TweetID: tweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "いいねの削除操作を完了できませんでした"})
			return
		}
		currentLikeStatus = false
	} else {
		_, err := tc.Queries.CreateLike(c.Request.Context(), db.CreateLikeParams{
			UserID:  userId,
			TweetID: tweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "いいねの登録操作を完了できませんでした"})
			return
		}
		currentLikeStatus = true
	}

	likeCount, err := tc.Queries.GetLikeCountByTweetID(c.Request.Context(), tweetId)

	touchActionResultRes := TouchActionResultResponse{
		TweetID:   tweetId,
		LikeCount: likeCount,
		IsLiked:   currentLikeStatus,
	}

	c.JSON(http.StatusOK, touchActionResultRes)

}
