package auth

import (
	"fmt"
	"github.com/iwariku/golang_twitter/utils"
	"net/smtp"
	"os"
)

type MailerInterface interface {
	SendActivationEmail(toMail string, token string) error
}

// docker環境用
type DevMailer struct {
	BaseURL     string
	SMTPAddress string
	FromEmail   string
}

// NewDevMailerは設定をするだけ
func NewDevMailer() *DevMailer {
	return &DevMailer{
		BaseURL:     utils.GetEnvOrDefault("APP_URL", "http://localhost:8080"),
		SMTPAddress: utils.GetEnvOrDefault("SMTP_ADDRESS", "mailcatcher:1025"),
		FromEmail:   utils.GetEnvOrDefault("MAIL_FROM", "noreply@example.com"),
	}
}

func (m *DevMailer) SendActivationEmail(toMail string, token string) error {
	activationURL := fmt.Sprintf("%v/activate?token=%s", m.BaseURL, token)

	subject := "Subject: アカウントを有効化してください\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><p>以下のリンクをクリックして、登録を完了させてください：</p><a href='%s'>%s</a></body></html>", activationURL, activationURL)
	msg := []byte(subject + mime + body)

	// SMTPを使って送信する
	return smtp.SendMail(m.SMTPAddress, nil, m.FromEmail, []string{toMail}, msg)
}

// 本番環境用
type ProdMailer struct {
	BaseURL     string
	SMTPAddress string
	FromEmail   string
	GooglePass  smtp.Auth //Gmailを使うため
}

// NewMailerは設定をするだけ
func NewProdMailer() *ProdMailer {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("GMAIL_USER"),
		os.Getenv("GMAIL_PASS"),
		"smtp.gmail.com",
	)

	return &ProdMailer{
		BaseURL:     utils.GetEnvOrDefault("FRONT_APP_URL", "http://localhost:3000"),
		SMTPAddress: os.Getenv("SMTP_ADDRESS"),
		FromEmail:   os.Getenv("GMAIL_USER"),
		GooglePass:  auth,
	}

}

func (m *ProdMailer) SendActivationEmail(toMail string, token string) error {
	activationURL := fmt.Sprintf("%v/activate?token=%s", m.BaseURL, token)

	subject := "Subject: アカウントを有効化してください\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(
		"<html><body><p>以下のリンクをクリックして、登録を完了させてください：</p><a href='%s'>%s</a></body></html>",
		activationURL, activationURL,
	)

	msg := []byte(subject + mime + body)

	return smtp.SendMail(
		m.SMTPAddress, m.GooglePass, m.FromEmail, []string{toMail}, msg)
}
