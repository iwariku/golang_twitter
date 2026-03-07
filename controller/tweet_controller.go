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

	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	tweet, err := tc.Queries.CreateTweet(c.Request.Context(), db.CreateTweetParams{
		UserID:  loggedUserId,
		Content: req.Content,
	})
	if err != nil {
		log.Printf("DBへの保存に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBへの保存に失敗しました"})
		return
	}

	tweetRes := TweetResponse{
		ID:      tweet.ID,
		UserID:  tweet.UserID,
		Content: tweet.Content,
	}

	c.JSON(http.StatusCreated, tweetRes)
}

// ===================
// ツイート一覧機能
// ===================
func (tc *TweetController) GetTweets(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
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

	dbTweets, err := tc.Queries.GetTweets(c.Request.Context(), db.GetTweetsParams{
		UserID: loggedUserId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBからの取得に失敗しました"})
		return
	}

	tweetRes := make([]TweetResponse, len(dbTweets))
	for i, t := range dbTweets {
		tweetRes[i] = TweetResponse{
			ID:           t.ID,
			UserID:       t.UserID,
			Content:      t.Content,
			LikeCount:    t.LikeCount,
			IsLiked:      t.IsLiked,
			RetweetCount: t.RetweetCount,
			IsRetweeted:  t.IsRetweeted,
			IsBookmarked: t.IsBookmarked,
		}
	}

	paginatedTweetsRes := PaginatedTweetsResponse{
		Tweets: tweetRes,
		Limit:  int(limit),
		Offset: int(offset),
		Count:  int(totalCount),
	}

	c.JSON(http.StatusOK, paginatedTweetsRes)
}

// ===================
// ツイート詳細機能
// ===================
func (tc *TweetController) GetTweet(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetTweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		log.Printf("パラメータ解析に失敗しました: %v", err)
		c.JSON(http.StatusBadRequest, "不正なリクエストです")
		return
	}

	dbTweet, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     targetTweetId,
	})
	if err != nil {
		log.Printf("DBからの取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, "DBからの取得に失敗しました")
		return
	}

	tweetRes := TweetResponse{
		ID:           dbTweet.ID,
		UserID:       dbTweet.UserID,
		Content:      dbTweet.Content,
		LikeCount:    dbTweet.LikeCount,
		IsLiked:      dbTweet.IsLiked,
		RetweetCount: dbTweet.RetweetCount,
		IsRetweeted:  dbTweet.IsRetweeted,
		IsBookmarked: dbTweet.IsBookmarked,
	}

	c.JSON(http.StatusOK, tweetRes)
}

// いいね機能
// バックエンドでツイートの有無の状態を確認する

// 1. 今何のツイートとユーザーIDなのかを確認する
// 2. 1.の情報を使いDBに該当するツイートがあるかを確認する
// 3. 2.の結果を元に条件分岐でcreateLikeかdeleteLikeかを決める
// 4. データを整形してレスポンスを返す
func (tc *TweetController) DeleteLike(c *gin.Context) {

	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	tweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	// Likeを持つ == DBにレコードがある
	hasLiked, err := tc.Queries.GetLikeExists(c.Request.Context(), db.GetLikeExistsParams{
		UserID:  loggedUserId,
		TweetID: tweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "条件に合致するツイートがDBにありません"})
		return
	}

	if hasLiked {
		err := tc.Queries.DeleteLike(c.Request.Context(), db.DeleteLikeParams{
			UserID:  loggedUserId,
			TweetID: tweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "いいねの削除操作を完了できませんでした"})
			return
		}
	}

	dbTweets, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     tweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同期に失敗しました"})
		return
	}

	touchActionResultRes := TouchActionResultResponse{
		TweetID:   tweetId,
		LikeCount: dbTweets.LikeCount,
		IsLiked:   dbTweets.IsLiked,
	}

	c.JSON(http.StatusOK, touchActionResultRes)

}

func (tc *TweetController) CreateLike(c *gin.Context) {

	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	tweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	// Likeを持つ == DBにレコードがある
	hasLiked, err := tc.Queries.GetLikeExists(c.Request.Context(), db.GetLikeExistsParams{
		UserID:  loggedUserId,
		TweetID: tweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "条件に合致するツイートがDBにありません"})
		return
	}

	if hasLiked == false {
		err := tc.Queries.CreateLike(c.Request.Context(), db.CreateLikeParams{
			UserID:  loggedUserId,
			TweetID: tweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "いいねの登録操作を完了できませんでした"})
			return
		}
	}

	dbTweets, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     tweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同期に失敗しました"})
		return
	}

	touchActionResultRes := TouchActionResultResponse{
		TweetID:   tweetId,
		LikeCount: dbTweets.LikeCount,
		IsLiked:   dbTweets.IsLiked,
	}

	c.JSON(http.StatusOK, touchActionResultRes)

}

// ツイート系にRetweetを入れる

// リツイートの登録、削除のロジックを作成する
// ログインしているユーザーの確認
// dbにレコードがあるか確認
// 条件分岐であるならdelete、あるならcreate
// いいねしたレスポンスを返す
// -> レスポンスの定義は一旦拡張せずに専用のレスポンス構造体を作成しようかな。いいねとリツイートはカウントがいるけど、ブックマークは要ら
func (tc *TweetController) DeleteRetweet(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetTweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	hasRetweeted, err := tc.Queries.GetRetweetExists(c.Request.Context(), db.GetRetweetExistsParams{
		UserID:  loggedUserId,
		TweetID: targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "条件に合致するツイートがありません"})
		return
	}

	if hasRetweeted {
		err := tc.Queries.DeleteRetweet(c.Request.Context(), db.DeleteRetweetParams{
			UserID:  loggedUserId,
			TweetID: targetTweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "リツイートの削除に失敗しました"})
			return
		}
	}

	dbTweet, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}

	touchActionRetweetRes := TouchActionRetweetResponse{
		TweetID:      targetTweetId,
		RetweetCount: dbTweet.RetweetCount,
		IsRetweeted:  dbTweet.IsRetweeted,
	}

	c.JSON(http.StatusOK, touchActionRetweetRes)
}

func (tc *TweetController) CreateRetweet(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetTweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	hasRetweeted, err := tc.Queries.GetRetweetExists(c.Request.Context(), db.GetRetweetExistsParams{
		UserID:  loggedUserId,
		TweetID: targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "条件に合致するツイートがありません"})
		return
	}

	if hasRetweeted == false {
		err := tc.Queries.CreateRetweet(c.Request.Context(), db.CreateRetweetParams{
			UserID:  loggedUserId,
			TweetID: targetTweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "リツイートの登録に失敗しました"})
			return
		}
	}

	dbTweet, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}

	touchActionRetweetRes := TouchActionRetweetResponse{
		TweetID:      targetTweetId,
		RetweetCount: dbTweet.RetweetCount,
		IsRetweeted:  dbTweet.IsRetweeted,
	}

	c.JSON(http.StatusOK, touchActionRetweetRes)
}

// ブックマーク削除
func (tc *TweetController) DeleteBookmark(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetTweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	hasBookmarked, err := tc.Queries.GetBookmarkExists(c.Request.Context(), db.GetBookmarkExistsParams{
		UserID:  loggedUserId,
		TweetID: targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "条件に合致するツイートがありません"})
		return
	}

	if hasBookmarked {
		err := tc.Queries.DeleteBookmark(c.Request.Context(), db.DeleteBookmarkParams{
			UserID:  loggedUserId,
			TweetID: targetTweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ブックマークの削除に失敗しました"})
			return
		}
	}

	dbTweet, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}

	touchActionBookmarkRes := TouchActionBookmarkResponse{
		TweetID:      targetTweetId,
		IsBookmarked: dbTweet.IsBookmarked,
	}

	c.JSON(http.StatusOK, touchActionBookmarkRes)
}

// ブックマーク登録
func (tc *TweetController) CreateBookmark(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetTweetId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tweet_idの形式が正しくありません"})
		return
	}

	hasBookmarked, err := tc.Queries.GetBookmarkExists(c.Request.Context(), db.GetBookmarkExistsParams{
		UserID:  loggedUserId,
		TweetID: targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "条件に合致するツイートがありません"})
		return
	}

	if hasBookmarked == false {
		err := tc.Queries.CreateBookmark(c.Request.Context(), db.CreateBookmarkParams{
			UserID:  loggedUserId,
			TweetID: targetTweetId,
		})
		if err != nil {
			log.Printf("データの更新に失敗しました")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "ブックマークの登録に失敗しました"})
			return
		}
	}

	dbTweet, err := tc.Queries.GetTweet(c.Request.Context(), db.GetTweetParams{
		UserID: loggedUserId,
		ID:     targetTweetId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}

	touchActionBookmarkRes := TouchActionBookmarkResponse{
		TweetID:      targetTweetId,
		IsBookmarked: dbTweet.IsBookmarked,
	}

	c.JSON(http.StatusOK, touchActionBookmarkRes)
}

// 選択したユーザーがどのツイートをリツイートしているのかを見る
// 必要なデータ。選択されたユーザーID、ログインしているユーザーID、LIMIT、OFFSET
// 使う関数: GetRetweetedTweetsByUserID
// api/retweetedtweet/:id?limit=1&offset=10

// 参考になるもの: ツイート一覧のページネーションのロジック、ユーザー詳細のツイート一覧、JSのfetch
func (tc *TweetController) GetRetweetedTweetsByUserID(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetUserId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "idの形式が違います"})
		return
	}

	limit, err := utils.ParseQueryInt32WithDefault(c, "limit", 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limitの形式が正しくありません"})
		return
	}

	offset, err := utils.ParseQueryInt32WithDefault(c, "offset", 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offsetの形式が正しくありません"})
		return
	}

	totalCount, err := tc.Queries.GetRetweetCountByUserID(c.Request.Context(), targetUserId)
	if err != nil {
		log.Printf("件数の取得に失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "件数取得に失敗しました"})
		return
	}

	dbTweet, err := tc.Queries.GetRetweetedTweetsByUserID(c.Request.Context(), db.GetRetweetedTweetsByUserIDParams{
		LoggedUserID: loggedUserId,
		TargetUserID: targetUserId,
		LimitVal:     limit,
		OffsetVal:    offset,
	})
	if err != nil {
		log.Printf("データの取得に失敗しました")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "データの取得に失敗しました"})
		return
	}

	tweetRes := make([]TweetResponse, len(dbTweet))
	for i, t := range dbTweet {
		tweetRes[i] = TweetResponse{
			ID:           t.ID,
			UserID:       t.UserID,
			Content:      t.Content,
			LikeCount:    t.LikeCount,
			IsLiked:      t.IsLiked,
			RetweetCount: t.RetweetCount,
			IsRetweeted:  t.IsRetweeted,
			IsBookmarked: t.IsBookmarked,
		}
	}

	paginatedTweetsRes := PaginatedTweetsResponse{
		Tweets: tweetRes,
		Limit:  int(limit),
		Offset: int(offset),
		Count:  int(totalCount),
	}

	c.JSON(http.StatusOK, paginatedTweetsRes)
}

// ログインしているユーザーのブックマークしたツイート一覧
func (tc *TweetController) GetBookmarkedTweetsByUserID(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "limitの形式が違います"})
		return
	}

	dbTweets, err := tc.Queries.GetBookmarkedTweetsByUserID(c.Request.Context(), db.GetBookmarkedTweetsByUserIDParams{
		UserID: loggedUserId,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBからの取得に失敗しました"})
		return
	}

	tweetRes := make([]TweetResponse, len(dbTweets))
	for i, t := range dbTweets {
		tweetRes[i] = TweetResponse{
			ID:           t.ID,
			UserID:       t.UserID,
			Content:      t.Content,
			LikeCount:    t.LikeCount,
			IsLiked:      t.IsLiked,
			RetweetCount: t.RetweetCount,
			IsRetweeted:  t.IsRetweeted,
			IsBookmarked: t.IsBookmarked,
		}
	}

	type BookmarkedTweetResponse struct {
		Tweets []TweetResponse `json:"tweets"`
	}

	bookmarkedTweetRes := BookmarkedTweetResponse{
		Tweets: tweetRes,
	}

	c.JSON(http.StatusOK, bookmarkedTweetRes)

}
