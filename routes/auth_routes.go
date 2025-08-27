package routes

import (
	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/AyKrimino/ScribeHost/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(authGroup *gin.RouterGroup, authController controller.AuthController) {
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/refresh", authController.RefreshToken)
	authGroup.POST("/logout", middleware.JwtAuthMiddleware(), authController.Logout)
	authGroup.POST("/verify-otp", authController.VerifyOTP)
	authGroup.POST("/resend-otp", authController.ResendOTP)
}
