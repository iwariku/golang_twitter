package controller

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TweetRequest struct {
	Content string `json:"content" binding:"required,max=140"`
}
