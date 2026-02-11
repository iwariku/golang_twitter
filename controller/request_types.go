package controller

type SignUpRequest struct {
	Mail     string `json:"mail" binding:"email"`
	Password string `json:"password"`
}

type TweetRequest struct {
	Content string `json:"content" binding:"required,max=140"`
}
