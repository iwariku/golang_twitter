package controller

import "golang_twitter/db"

// formatTweetResponseはDBモデルからAPIレスポンス用の構造体に変換する
func FormatTweetResponse(userID int32, content string) TweetResponse {
	return TweetResponse{
		UserID:  userID,
		Content: content,
	}
}

// ページネーション付きで返す
// ユーザー一覧(フォローフォロワーのタスク)でもページネーションは使うため先に分けておく
func FormatPaginatedTweetsResponse(tweets []db.Tweet, limit, offset int32, totalCount int64) PaginatedTweetsResponse {
	var tweetsRes []TweetResponse
	for _, t := range tweets {
		tweetsRes = append(tweetsRes, FormatTweetResponse(t.UserID, t.Content))
	}

	return PaginatedTweetsResponse{
		Tweets: tweetsRes,
		Limit:  int(limit),
		Offset: int(offset),
		Count:  int(totalCount),
	}
}
