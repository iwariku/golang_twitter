package auth

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendActivationEmail(toMail string, token string) error {
	// docker-compose.yamlから取得、ハードコードを避けるため
	baseURL := os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	activationURL := fmt.Sprintf("%v/activate?token=%s", baseURL, token)

	subject := "Subject: アカウントを有効化してください"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := fmt.Sprintf("<html><body><p>以下のリンクをクリックして、登録を完了させてください：</p><a href='%s'>%v</a></body></html>", activationURL, activationURL)
	msg := []byte(subject + mime + body)

	from := "noreply@example.com"
	to := []string{toMail}
	smtpAddress := os.Getenv("SMTP_ADDRESS")

	// SMTPを使って送信する
	err := smtp.SendMail(smtpAddress, nil, from, to, msg)
	if err != nil {
		return err
	}

	return nil
}
