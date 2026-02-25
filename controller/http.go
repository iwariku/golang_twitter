package controller

import "time"

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
}

// tweet
type TweetRequest struct {
	Content string `json:"content" binding:"required,max=140"`
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

type TouchActionResultResponse struct {
	TweetID   int32
	LikeCount int64
	IsLiked   bool
}
