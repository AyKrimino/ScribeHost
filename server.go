package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/AyKrimino/ScribeHost/middleware"
	"github.com/AyKrimino/ScribeHost/repository"
	"github.com/AyKrimino/ScribeHost/routes"
	"github.com/AyKrimino/ScribeHost/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	var err error

	setupLoggerOutput()

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	err = repository.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// repositories
	userRepo := repository.NewUserRepository()
	authRepo := repository.NewAuthRepository()
	refreshTokenRepo := repository.NewRefreshTokenRepository()

	// services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(authRepo, refreshTokenRepo)

	// controllers
	userController := controller.NewUserController(userService)
	authController := controller.NewAuthController(authService)

	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	routes.SetupRoutes(router, userController, authController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = router.Run(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupLoggerOutput() {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	logFile := &lumberjack.Logger{
		Filename:   "logs/application.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	gin.DefaultWriter = io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(io.MultiWriter(logFile, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
