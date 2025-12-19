package routes

import (
	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/AyKrimino/ScribeHost/middleware"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(authGroup *gin.RouterGroup, authController controller.AuthController) {
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
