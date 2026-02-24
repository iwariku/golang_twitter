package controller

import "golang_twitter/db"

// formatTweetResponseはDBモデルからAPIレスポンス用の構造体に変換する
// 拡張性を意識し、引数はDBモデルそのものを渡す。
func FormatTweetResponse(t db.Tweet) TweetResponse {
	return TweetResponse{
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
