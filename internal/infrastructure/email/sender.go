package email

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AyKrimino/ScribeHost/internal/infrastructure/config"
	"gopkg.in/gomail.v2"
)

func SendEmail(templateFile, placeholder, replacement, subject, email string) error {
	config.LoadEnv()

	host := os.Getenv("EMAIL_SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		return err
	}
	username := os.Getenv("EMAIL_SMTP_USER")
	password := os.Getenv("EMAIL_SMTP_PASSWORD")
	receiverEmail := os.Getenv("EMAIL_FROM")

	template, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	currYear, _, _ := time.Now().Date()

	htmlBody := strings.ReplaceAll(string(template), placeholder, replacement)
	htmlBody = strings.ReplaceAll(htmlBody, "%%YEAR%%", fmt.Sprint(currYear))

	m := gomail.NewMessage()
	m.SetHeader("From", receiverEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendOTPByEmail(email, otp string) error {
	return SendEmail(
		"templates/emails/verify_email.html",
		"%%OTP%%",
		otp,
		"Email Verification - ScribeHost",
		email,
	)
}

func SendPasswordResetEmail(email, resetLink string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	resetLink = frontendURL + resetLink

	return SendEmail(
		"templates/emails/password_reset_email.html",
		"%%RESET_LINK%%",
		resetLink,
		"Password Reset - ScribeHost",
		email,
	)
}
