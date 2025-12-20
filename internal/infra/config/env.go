package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvType string

const (
	Dev  EnvType = "development"
	Prod EnvType = "production"
)

var CurrEnv EnvType

func LoadEnv() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = string(Dev)
	}

	switch env {
	case string(Prod):
		CurrEnv = Prod
	default:
		CurrEnv = Dev
	}

	envFile := ".env"
	if CurrEnv == Prod {
		envFile += ".production"
	} else {
		envFile += ".local"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(".env file not found, relying on environment variables.")
		} else {
			log.Printf("Failed to load .env file: %v", err)
		}
	} else {
		log.Printf("Environment loaded: %s (using %s)", CurrEnv, envFile)
	}
}
