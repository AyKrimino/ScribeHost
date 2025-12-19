package server

import (
	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/AyKrimino/ScribeHost/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userController controller.UserController, authController controller.AuthController) {
	v1 := router.Group("api/v1")
	{
		authGroup := v1.Group("/auth")
		setupAuthRoutes(authGroup, authController)

		userGroup := v1.Group("/users")
		setupUserRoutes(userGroup, userController)
	}
}

func setupAuthRoutes(authGroup *gin.RouterGroup, authController controller.AuthController) {
	authGroup.POST(
		"/register",
		middleware.RateLimiter("register"),
		authController.Register,
	)

	authGroup.POST(
		"/login",
		middleware.RateLimiter("login"),
		authController.Login,
	)

	authGroup.POST(
		"/refresh",
		authController.RefreshToken,
	)

	authGroup.POST(
		"/logout",
		middleware.JwtAuthMiddleware(),
		authController.Logout,
	)

	authGroup.POST(
		"/verify-otp",
		middleware.RateLimiter("verify-otp"),
		authController.VerifyOTP,
	)

	authGroup.POST(
		"/resend-otp",
		middleware.RateLimiter("resend-otp"),
		authController.ResendOTP,
	)

	authGroup.POST(
		"/forgot-password",
		middleware.RateLimiter("forgot-password"),
		authController.ForgotPassword,
	)

	authGroup.POST(
		"/reset-password",
		middleware.RateLimiter("reset-password"),
		authController.ResetPassword,
	)
}

func setupUserRoutes(userGroup *gin.RouterGroup, userController controller.UserController) {
	userGroup.POST("", userController.CreateUser)
}
