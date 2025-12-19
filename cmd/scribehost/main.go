package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/AyKrimino/ScribeHost/helper"
	"github.com/AyKrimino/ScribeHost/internal/auth"
	"github.com/AyKrimino/ScribeHost/internal/infrastructure/database"
	"github.com/AyKrimino/ScribeHost/internal/server"
	"github.com/AyKrimino/ScribeHost/middleware"
	"github.com/AyKrimino/ScribeHost/repository"
	"github.com/AyKrimino/ScribeHost/service"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	var err error

	setupLoggerOutput()

	helper.LoadEnv()

	db, err := database.InitMySQLDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	ctx := context.Background()

	redisClient, err := helper.CreateRedisClientConn(ctx)
	if err != nil {
		log.Fatalf("Failed to create redis client: %v", err)
	}
	defer redisClient.Close()

	middleware.SetRedisClient(redisClient)

	// repositories
	userRepo := repository.NewUserRepository()
	authRepo := auth.NewAuthRepository(db)
	refreshTokenRepo := repository.NewRefreshTokenRepository()
	otpRedisRepo := auth.NewOtpRedisRepo(redisClient, ctx)
	passwordResetRedisRepo := repository.NewPasswordResetRedisRepo(redisClient, ctx)

	// services
	userService := service.NewUserService(userRepo)
	authService := auth.NewAuthService(
		authRepo,
		refreshTokenRepo,
		otpRedisRepo,
		passwordResetRedisRepo,
	)

	// controllers
	userController := controller.NewUserController(userService)
	authController := auth.NewAuthController(authService)

	router := gin.New()
	router.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	server.SetupRoutes(router, userController, authController)

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
