package repository

import (
	"fmt"
	"os"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDB() error {
	var err error

	err = godotenv.Load()
	if err != nil {
		return fmt.Errorf("failed to load .env file: %w", err)
	}

	var (
		databaseUsername = os.Getenv("DB_USERNAME")
		databasePassword = os.Getenv("DB_PASSWORD")
		databaseHost     = os.Getenv("DB_HOST")
		databasePort     = os.Getenv("DB_PORT")
		databaseName     = os.Getenv("DB_NAME")
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		databaseUsername,
		databasePassword,
		databaseHost,
		databasePort,
		databaseName,
	)

	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	dbResult := DB.AutoMigrate(
		&entity.User{},
		&entity.Profile{},
		&entity.RefreshToken{},
	)
	if dbResult.Error != nil {
		return fmt.Errorf("failed to migrate database: %w", dbResult.Error)
	}

	return nil
}
