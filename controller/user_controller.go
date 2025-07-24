package controller

import (
	"net/http"

	"github.com/AyKrimino/ScribeHost/dto"
	"github.com/AyKrimino/ScribeHost/service"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
}

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (c *userController) CreateUser(ctx *gin.Context) {
	var (
		req dto.CreateUserRequestDto
		res *dto.UserResponseDto
		err error
	)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	res, err = c.userService.CreateUser(req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
