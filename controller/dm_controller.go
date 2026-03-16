package controller

import (
	"golang_twitter/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type DmController struct {
	Queries *db.Queries
	Redis   *redis.Client
}

// リクエストする時に必要な情報
type GroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type GroupResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

// リクエストはグループ名のみ
func (dc *DmController) CreateGroup(c *gin.Context) {
	var req GroupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSONの形式が違います: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSONの形式が違います"})
		return
	}

	groupName, err := dc.Queries.CreateGroup(c.Request.Context(), req.Name)
	if err != nil {
		log.Printf("DBへの保存に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBへの保存に失敗しました"})
		return
	}

	groupNameRes := GroupResponse{
		ID:   groupName.ID,
		Name: groupName.Name,
	}

	c.JSON(http.StatusOK, groupNameRes)
}
