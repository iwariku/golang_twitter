package controller

type SignUpRequest struct {
	Mail     string `json:"mail" binding:"email"`
	Password string `json:"password"`
}
