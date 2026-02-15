package auth

import (
	"fmt"
	"golang_twitter/utils"
	"net/smtp"
)

type MailerInterface interface {
	SendActivationEmail(toMail string, token string) error
}

type Mailer struct {
	BaseURL     string
	SMTPAddress string
	FromEmail   string
}

// NewMailerは設定をするだけ
func NewMailer() *Mailer {
	return &Mailer{
		BaseURL:     utils.GetEnvOrDefault("APP_URL", "http://localhost:8080"),
		SMTPAddress: utils.GetEnvOrDefault("SMTP_ADDRESS", "mailcatcher:1025"),
		FromEmail:   utils.GetEnvOrDefault("MAIL_FROM", "noreply@example.com"),
	}
}

func (m *Mailer) SendActivationEmail(toMail string, token string) error {
	activationURL := fmt.Sprintf("%v/activate?token=%s", m.BaseURL, token)

	subject := "Subject: アカウントを有効化してください\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><p>以下のリンクをクリックして、登録を完了させてください：</p><a href='%s'>%s</a></body></html>", activationURL, activationURL)
	msg := []byte(subject + mime + body)

	// SMTPを使って送信する
	return smtp.SendMail(m.SMTPAddress, nil, m.FromEmail, []string{toMail}, msg)
}
