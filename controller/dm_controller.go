package controller

import (
	"fmt"
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

type MessageResponse struct {
	ID      int32  `json:"id"`
	UserID  int32  `json:"user_id"`
	Message string `json:"message"`
}

// リクエストはグループ名のみ
func (dc *DmController) CreateGroup(c *gin.Context) {
	type GroupRequest struct {
		Name string `json:"name" binding:"required"`
	}

	type GroupResponse struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}

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

func (dc *DmController) AddMemberToGroup(c *gin.Context) {
	type AddmemberRequest struct {
		UserID  int32 `json:"user_id"`
		GroupID int32 `json:"group_id"`
	}

	type GroupMemberResponse struct {
		UserID    int32 `json:"user_id"`
		DmGroupID int32 `json:"dm_group_id"`
	}

	var req AddmemberRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("JSONの形式が違います: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSONの形式が違います"})
		return
	}

	// ユーザーがグループにすでに追加されているかどうかを確認する
	hasDmGroup, err := dc.Queries.AlreadyAddUserToGroup(c.Request.Context(), db.AlreadyAddUserToGroupParams{
		UserID:    req.UserID,
		DmGroupID: req.GroupID,
	})

	if hasDmGroup {
		log.Printf("ユーザー(ID:%d)は既にグループ(ID:%d)に存在します", req.UserID, req.GroupID)
		c.JSON(http.StatusConflict, gin.H{"error": "このユーザーは既にグループに追加されています"})
		return
	}
	fmt.Println("入力されたユーザーを追加可能です", hasDmGroup)

	groupMember, err := dc.Queries.AddMemberToGroup(c.Request.Context(), db.AddMemberToGroupParams{
		UserID:    req.UserID,
		DmGroupID: req.GroupID,
	})
	if err != nil {
		log.Printf("DBへの保存に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DBへの保存に失敗しました"})
		return
	}

	groupMemberRes := GroupMemberResponse{
		UserID:    groupMember.UserID,
		DmGroupID: groupMember.DmGroupID,
	}

	c.JSON(http.StatusOK, groupMemberRes)

}

// 必要なデータ、誰が: user_id,どこに: dm_group_id、何を: message
func (dc *DmController) CreateMessage(c *gin.Context) {
	type MessageRequest struct {
		Message string `json:"message" binding:"required"`
	}

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
		ID:      Message.ID,
		UserID:  Message.UserID,
		Message: Message.Message,
	}

	c.JSON(http.StatusOK, messageRes)

}

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
			ID:      m.ID,
			UserID:  m.UserID,
			Message: m.Message,
		}
	}

	type MessagesResponse struct {
		Messages []MessageResponse `json:"messages"`
	}

	messageListRes := MessagesResponse{
		Messages: messagesRes,
	}

	c.JSON(http.StatusOK, messageListRes)
}

// 構造体はこのメソッド内でいいのか。ここだけであれば問題がない
// dbから取得したデータとレスポンスの変数名の命名規則を明確にする必要がある
func (dc *DmController) GetGroups(c *gin.Context) {
	loggedUserId, err := GetUserIDFromContext(c)
	if err != nil {
		log.Printf("ログインチェックの失敗: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインが必要です"})
		return
	}

	dbGroups, err := dc.Queries.GetGroups(c.Request.Context(), loggedUserId)
	if err != nil {
		log.Printf("グループの取得に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "グループの取得に失敗しました"})
		return
	}

	type GroupResponse struct {
		ID   int32  `json:"id"`
		Name string `json:"name"`
	}

	groupsRes := make([]GroupResponse, len(dbGroups))
	for i, g := range dbGroups {
		groupsRes[i] = GroupResponse{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	type GroupsResponse struct {
		Groups []GroupResponse `json:"groups"`
	}

	groupListRes := GroupsResponse{
		Groups: groupsRes,
	}

	c.JSON(http.StatusOK, groupListRes)

}
