package routes

import (
	"github.com/AyKrimino/ScribeHost/controller"
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(userGroup *gin.RouterGroup, userController controller.UserController) {
	userGroup.POST("", userController.CreateUser)
}
