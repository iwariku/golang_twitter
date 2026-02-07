package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func (uc *UserController) Activate(c *gin.Context) {
	// URLから"token"という名前のパラメータを取得する
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "トークンが必要です"})
		return
	}

	err := uc.Queries.ActivateUser(c.Request.Context(), pgtype.Text{String: token, Valid: true})
	if err != nil {
		log.Printf("有効化エラー: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "有効化に失敗しました"})
	}

	c.JSON(http.StatusOK, gin.H{"message": "アカウントが有効になりました。ログインが可能な状態です。"})
}
