package controller

import (
	"time"
)

// リクエストレスポンスの構造体を全てここにまとめる

// user
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserResponse struct {
	UserName         string    `json:"user_name"`
	SelfIntroduction string    `json:"self_introduction"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	ProfileImage     string    `json:"profile_image"`
	FollowingCount   int64     `json:"following_count"`
	FollowerCount    int64     `json:"follower_count"`
	IsFollowed       bool      `json:"is_followed"`
}

// tweet
type TweetRequest struct {
	Content string `json:"content" binding:"required,max=140"`
}

// 構造体の中に構造体を入れる
type TweetResponse struct {
	ID           int32  `json:"id"`
	UserID       int32  `json:"user_id"`
	Content      string `json:"content"`
	LikeCount    int64  `json:"like_count"`
	IsLiked      bool   `json:"is_liked"`
	RetweetCount int64  `json:"retweet_count"`
	IsRetweeted  bool   `json:"is_retweeted"`
	IsBookmarked bool   `json:"is_bookmarked"`
}

// この構造体はそのままでいいんじゃない？ツイート一覧を返却するのであれば
type PaginatedTweetsResponse struct {
	Tweets []TweetResponse `json:"tweets"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Count  int             `json:"count"`
}

type BookmarkedTweetResponse struct {
	Tweets []TweetResponse `json:"tweets"`
}

type TouchActionResultResponse struct {
	TweetID   int32 `json:"tweet_id"`
	LikeCount int64 `json:"like_count"`
	IsLiked   bool  `json:"is_liked"`
}

type TouchActionRetweetResponse struct {
	TweetID      int32 `json:"tweet_id"`
	RetweetCount int64 `json:"retweet_count"`
	IsRetweeted  bool  `json:"is_retweeted"`
}

type TouchActionBookmarkResponse struct {
	TweetID      int32 `json:"tweet_id"`
	IsBookmarked bool  `json:"is_bookmarked"`
}

type FollowResponse struct {
	FollowerID  int32 `json:"follower_id"`  // 誰が
	FollowingID int32 `json:"following_id"` // 誰を
	IsFollowed  bool  `json:"is_followed"`  // フォローしているか？true/false
}
