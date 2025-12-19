package database

import (
	"fmt"
	"os"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/AyKrimino/ScribeHost/helper"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitMySQLDB() (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	helper.LoadEnv()

	var (
		databaseUsername = os.Getenv("DB_USERNAME")
		databasePassword = os.Getenv("DB_PASSWORD")
		databaseHost     = os.Getenv("DB_HOST")
		databasePort     = os.Getenv("DB_PORT")
		databaseName     = os.Getenv("DB_NAME")
	)

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local",
		databaseUsername,
		databasePassword,
		databaseHost,
		databasePort,
		databaseName,
	)

	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	dbResult := db.AutoMigrate(
		&entity.User{},
		&entity.Profile{},
		&entity.RefreshToken{},
	)
	if dbResult.Error != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", dbResult.Error)
	}

	return db, nil
}
