package controller

import (
	"golang_twitter/db"
	"golang_twitter/utils"
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

type GroupMemberResponse struct {
	UserID    int32 `json:"user_id"`
	DmGroupID int32 `json:"dm_group_id"`
}

// メッセージを送信する際に必要な情報
type MessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type MessageResponse struct {
	UserID  int32  `json:"user_id"`
	Message string `json:"message"`
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

// あるユーザーが作成したグループに、自分と相手が追加される(単一責任の原則について考慮する)
// トランザクションを使って一連の動きにするのであればログインユーザーを使う
// 別のユーザーを追加するときは別のメソッドを定義する方がいいと思う(単一責任の原則とRESTの設計に準ずる)
// グループに、自分と相手が追加される
func (dc *DmController) AddMemberToGroup(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetGroupId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		log.Printf("パラメータ解析に失敗しました: %v", err)
		c.JSON(http.StatusBadRequest, "idの形式が違います")
		return
	}

	groupMember, err := dc.Queries.AddMemberToGroup(c.Request.Context(), db.AddMemberToGroupParams{
		UserID:    loggedUserId,
		DmGroupID: targetGroupId,
	})

	groupMemberRes := GroupMemberResponse{
		UserID:    loggedUserId,
		DmGroupID: groupMember.DmGroupID,
	}

	c.JSON(http.StatusOK, groupMemberRes)

}

// 必要なデータ、誰が: user_id,どこに: dm_group_id、何を: message
func (dc *DmController) CreateMessage(c *gin.Context) {
	var req MessageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSONの形式が違います: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSONの形式が違います"})
		return
	}

	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	targetGroupId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		log.Printf("パラメータ解析に失敗しました: %v", err)
		c.JSON(http.StatusBadRequest, "idの形式が違います")
		return
	}

	Message, err := dc.Queries.CreateMessage(c.Request.Context(), db.CreateMessageParams{
		UserID:    loggedUserId,
		DmGroupID: targetGroupId,
		Message:   req.Message,
	})
	if err != nil {
		log.Printf("メッセージの作成に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの作成に失敗しました"})
		return
	}

	messageRes := MessageResponse{
		UserID:  Message.UserID,
		Message: Message.Message,
	}

	c.JSON(http.StatusOK, messageRes)

}

// user_idはどうする？今のままだと自分と違うグループのユーザーの見れてしまう
// SQLのWHERE句にuser_idもつけるのが一番簡単。ただ、ユーザー数が多くいる場合はどうする？
func (dc *DmController) GetMessagesByGroupID(c *gin.Context) {
	targetGroupId, err := utils.ParseParamInt32(c, "id")
	if err != nil {
		log.Printf("パラメータ解析に失敗しました: %v", err)
		c.JSON(http.StatusBadRequest, "idの形式が違います")
		return
	}

	dbMessages, err := dc.Queries.GetMessagesByGroupID(c.Request.Context(), targetGroupId)
	if err != nil {
		log.Printf("メッセージの取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの取得に失敗しました"})
		return
	}

	messagesRes := make([]MessageResponse, len(dbMessages))
	for i, m := range dbMessages {
		messagesRes[i] = MessageResponse{
			UserID:  m.UserID,
			Message: m.Message,
		}
	}

	type MessagesResponse struct {
		Messages []MessageResponse `json:"messages"`
	}

	filtedMessagesRes := MessagesResponse{
		Messages: messagesRes,
	}

	c.JSON(http.StatusOK, filtedMessagesRes)
}
