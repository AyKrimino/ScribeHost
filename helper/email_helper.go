package helper

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendOTPByEmail(email, otp string) error {
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found, relying on environment variables.")
		} else {
			log.Printf("Failed to load .env file: %v", err)
		}
	} else {
		log.Println("Successfully loaded .env file.")
	}

	host := os.Getenv("EMAIL_SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))
	if err != nil {
		return err
	}
	username := os.Getenv("EMAIL_SMTP_USER")
	password := os.Getenv("EMAIL_SMTP_PASSWORD")
	receiverEmail := os.Getenv("EMAIL_FROM")

	m := gomail.NewMessage()
	m.SetHeader("From", receiverEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", fmt.Sprintf("Hello <b>User</b> this is your otp <i>%s</i>", otp))

	d := gomail.NewDialer(host, port, username, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
