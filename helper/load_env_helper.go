package helper

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
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
}
