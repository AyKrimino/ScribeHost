package routes

import (
	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/gin-gonic/gin"
)

func SetupAuthRoutes(authGroup *gin.RouterGroup, authController controller.AuthController) {
	authGroup.POST("/register", authController.Register)
	authGroup.POST("/login", authController.Login)
	authGroup.POST("/refresh", authController.RefreshToken)
}
