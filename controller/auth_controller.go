package controller

import (
	"log"
	"net/http"

	"github.com/AyKrimino/ScribeHost/dto"
	"github.com/AyKrimino/ScribeHost/errors"
	"github.com/AyKrimino/ScribeHost/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type authController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var (
		req dto.RegisterRequestDto
		res *dto.RegisterResponseDto
		err error
	)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("Register binding error: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": "The request data is invalid. Please check your input and try again.",
		})
		return
	}

	res, err = c.authService.Register(req)
	if err != nil {
		if errors.IsObjectAlreadyExists(err) {
			log.Printf("Registration conflict: %v", err)

			ctx.JSON(http.StatusConflict, gin.H{
				"error":   "Registration failed",
				"message": "A user with this email already exists. Please use a different email or try logging in.",
			})
			return
		}

		log.Printf("Registration error for %s: %v", req.Email, err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"message": "An internal error occurred. Please try again later.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *authController) Login(ctx *gin.Context) {
}

func (c *authController) RefreshToken(ctx *gin.Context) {
}
