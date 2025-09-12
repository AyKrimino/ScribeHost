package helper

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

func SendOTPByEmail(email, otp string) error {
	LoadEnv()

	host := os.Getenv("EMAIL_SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		return err
	}
	username := os.Getenv("EMAIL_SMTP_USER")
	password := os.Getenv("EMAIL_SMTP_PASSWORD")
	receiverEmail := os.Getenv("EMAIL_FROM")

	template, err := os.ReadFile("templates/emails/verify_email.html")
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	htmlBody := strings.Replace(string(template), "%%OTP%%", otp, -1)

	m := gomail.NewMessage()
	m.SetHeader("From", receiverEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification - ScribeHost")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendPasswordResetEmail(email, resetLink string) error {
	LoadEnv()

	host := os.Getenv("EMAIL_SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		return err
	}
	username := os.Getenv("EMAIL_SMTP_USER")
	password := os.Getenv("EMAIL_SMTP_PASSWORD")
	receiverEmail := os.Getenv("EMAIL_FROM")

	template, err := os.ReadFile("templates/emails/password_reset_email.html")
	if err != nil {
		return fmt.Errorf("failed to read email template: %w", err)
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	resetLink = frontendURL + resetLink
	htmlBody := strings.Replace(string(template), "%%RESET_LINK%%", resetLink, -1)

	m := gomail.NewMessage()
	m.SetHeader("From", receiverEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset - ScribeHost")
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
