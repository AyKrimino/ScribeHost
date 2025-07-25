package routes

import (
	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userController controller.UserController, authController controller.AuthController) {
	v1 := router.Group("api/v1")
	{
		authGroup := v1.Group("/auth")
		SetupAuthRoutes(authGroup, authController)

		userGroup := v1.Group("/users")
		SetupUserRoutes(userGroup, userController)
	}
}
