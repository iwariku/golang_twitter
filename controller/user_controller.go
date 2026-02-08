package controller

import (
	"golang_twitter/db"
	"golang_twitter/services/auth"
	"golang_twitter/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

// このuser.goではuserに関する処理を持つファイル
// Signup
// Login
// Logout
// GetUserProfile
// UpdateUser

type UserController struct {
	Queries *db.Queries
	Mailer  auth.MailerInterface
}

// SignUpの流れ
// UserController型のポインタを示す変数がSingUpというメソッドを持つ
// SingUpメソッドは(c *gin.Context)を引数に取る。*gin.ContextはGinフレームワークがHTTPリクエストの時に自動的に作ってくれる
// reqというサインアップに必要なプロパティを持つ変数を宣言
// ShouldBindJSONで受け取ったJSONを使いreqを上書き。書き換える内容は、入力されたリクエスト(メールアドレスとパスワード)
// パスワードが何文字以上？大文字、小文字等の要件を満たしているかをチェック
// パスワードチェックに問題がなかったらハッシュ化
// DBに登録する(CreateUserメソッドはsqlcで自動的に作成されたもの)
func (uc *UserController) SignUp(c *gin.Context) {
	var req SignUpRequest

	// リクエスト情報などが詰まっている「c」からJSONを取り出して、reqという箱に詰め替える
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("メールアドレスの形式で入力してください: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "メールアドレスの形式で入力してください"})
		return
	}

	if err := validatePassword(req.Password); err != nil {
		log.Printf("バリデーションエラー: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateToken()
	if err != nil {
		log.Printf("トークン生成失敗: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー内部でエラーが発生"})
		return
	}

	user, err := uc.Queries.CreateUser(c.Request.Context(), db.CreateUserParams{
		Mail:            req.Mail,
		Password:        string(hashedPassword),
		IsActive:        pgtype.Bool{Bool: false, Valid: true},
		ActivationToken: pgtype.Text{String: token, Valid: true},
	})
	if err != nil {
		log.Printf("入力されたメールアドレスがすでに使用されています: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力されたメールアドレスがすでに使用されています"})
		return
	}

	err = uc.Mailer.SendActivationEmail(user.Mail, token)
	if err != nil {
		log.Printf("メールの送信に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "メールの送信に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

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
