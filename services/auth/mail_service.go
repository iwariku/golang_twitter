package auth

import (
	"fmt"
	"net/smtp"
	"os"
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
		BaseURL:     getEnvOrDefault("APP_URL", "http://localhost:8080"),
		SMTPAddress: getEnvOrDefault("SMTP_ADDRESS", "mailcatcher:1025"),
		FromEmail:   getEnvOrDefault("MAIL_FROM", "noreply@example.com"),
	}
}

// 環境変数を取得してif文で判定するコードが2つあったため、まとめて関数にした。NewMailerで使用
func getEnvOrDefault(key, fallback string) string {
	if result := os.Getenv(key); result != "" {
		return result
	}
	return fallback
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
