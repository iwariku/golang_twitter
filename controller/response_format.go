package controller

import "golang_twitter/db"

// ==================
// ツイート(いいね無し)のレスポンスを返却
// ==================

// formatTweetResponseはDBモデルからAPIレスポンス用の構造体に変換する
// 拡張性を意識し、引数はDBモデルそのものを渡す。
func FormatTweetResponse(t db.Tweet) TweetResponse {
	return TweetResponse{
		ID:      t.ID,
		UserID:  t.UserID,
		Content: t.Content,
	}
}

// ページネーション付きで返す
// ユーザー一覧(フォローフォロワーのタスク)でもページネーションは使うため先に分けておく
func FormatPaginatedTweetsResponse(tweets []db.Tweet, limit, offset int32, totalCount int64) PaginatedTweetsResponse {
	var tweetsRes []TweetResponse
	for _, t := range tweets {
		tweetsRes = append(tweetsRes, FormatTweetResponse(t))
	}

	return PaginatedTweetsResponse{
		Tweets: tweetsRes,
		Limit:  int(limit),
		Offset: int(offset),
		Count:  int(totalCount),
	}
}

// ==================
// ツイート(いいね数と状態を持つ)レスポンスを返却
// ==================

// いいね月ツイートをAPIレスポンス用の構造体に変換する
// リツイートとブックマーク機能もここに足す。その時に命名変更
// ツイートだけをレスポンスするFormatTweetResponseを使わずこちらをメインにする
func FormatTweetWithLikeResponse(t db.GetTweetsWithLikesRow) TweetResponse {
	return TweetResponse{
		ID:        t.ID,
		UserID:    t.UserID,
		Content:   t.Content,
		LikeCount: t.LikeCount,
		IsLiked:   t.IsLiked,
	}
}

// ページネーション付きで返す
// いいね付きなので、db.GetTweetsWithLikesRow型にする
// もしリツイートとブックマークが追加されたとしてもTweetResponseの構造体を変えるだけでいい
func FormatPaginatedWithLikeTweetsResponse(dbTweets []db.GetTweetsWithLikesRow, limit, offset int32, totalCount int64) PaginatedTweetsResponse {
	tweetsRes := make([]TweetResponse, len(dbTweets))
	for i, t := range dbTweets {
		tweetsRes[i] = FormatTweetWithLikeResponse(t) // 新設した関数を利用
	}

	return PaginatedTweetsResponse{
		Tweets: tweetsRes,
		Limit:  int(limit),
		Offset: int(offset),
		Count:  int(totalCount),
	}
}
