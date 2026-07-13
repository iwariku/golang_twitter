package auth

import (
	"fmt"
	"net"
	"net/smtp"
	"os"

	"github.com/iwariku/golang_twitter/utils"
)

type MailerInterface interface {
	SendActivationEmail(toMail string, token string) error
}

// docker環境用
type DevMailer struct {
	BaseURL   string
	SMTPHost  string
	SMTPPort  string
	FromEmail string
}

// NewDevMailerは設定をするだけ
func NewDevMailer() *DevMailer {
	return &DevMailer{
		BaseURL:   utils.GetEnvOrDefault("APP_URL", "http://localhost:8080"),
		SMTPHost:  utils.GetEnvOrDefault("SMTP_HOST", "mailcatcher"),
		SMTPPort:  utils.GetEnvOrDefault("SMTP_PORT", "1025"),
		FromEmail: utils.GetEnvOrDefault("FROM_EMAIL", "noreply@example.com"),
	}
}

func (m *DevMailer) SendActivationEmail(toMail string, token string) error {
	activationURL := fmt.Sprintf("%v/activate?token=%s", m.BaseURL, token)

	subject := "Subject: アカウントを有効化してください\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><p>以下のリンクをクリックして、登録を完了させてください：</p><a href='%s'>%s</a></body></html>", activationURL, activationURL)
	msg := []byte(subject + mime + body)

	// SMTPを使って送信する
	return smtp.SendMail(net.JoinHostPort(m.SMTPHost, m.SMTPPort), nil, m.FromEmail, []string{toMail}, msg)
}

// 本番環境用
type ProdMailer struct {
	BaseURL    string
	SMTPHost   string
	SMTPPort   string
	FromEmail  string
	GooglePass smtp.Auth //Gmailを使うため
}

// NewMailerは設定をするだけ
func NewProdMailer() *ProdMailer {

	auth := smtp.PlainAuth(
		"",
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"),
		os.Getenv("SMTP_HOST"),
	)

	return &ProdMailer{
		BaseURL:    utils.GetEnvOrDefault("FRONT_APP_URL", "http://localhost:3000"),
		SMTPHost:   os.Getenv("SMTP_HOST"),
		SMTPPort:   os.Getenv("SMTP_PORT"),
		FromEmail:  os.Getenv("FROM_EMAIL"),
		GooglePass: auth,
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
		net.JoinHostPort(m.SMTPHost, m.SMTPPort), m.GooglePass, m.FromEmail, []string{toMail}, msg)
}
