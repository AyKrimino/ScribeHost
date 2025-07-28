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
				"details": "A user with this email already exists. Please use a different email or try logging in.",
			})
			return
		}

		log.Printf("Registration error for %s: %v", req.Email, err)

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"details": "An internal error occurred. Please try again later.",
		})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *authController) Login(ctx *gin.Context) {
	var (
		req dto.LoginRequestDto
		res *dto.LoginResponseDto
		err error
	)

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Printf("Login binding error: %v", err)

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": "The request data is invalid. Please check your input and try again.",
		})
		return
	}

	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	res, err = c.authService.Login(req, userAgent, clientIP)
	if err != nil {
		log.Printf("Login error for %s: %v", req.Email, err)

		if errors.IsObjectNotFoundError(err) || errors.IsInvalidCredentialsError(err) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"details": "The email or password you entered is incorrect.",
			})
			return
		}

		if errors.IsObjectNotActiveError(err) {
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   "Account inactive",
				"details": "Your account is currently inactive. Please contact support.",
			})
			return
		}

		// TODO: check for email verified error type

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Login failed",
			"details": "An internal error occured. Please try again later.",
		})
		return
	}

	accessTokenMaxAge := int(res.AccessTokenTTL.Seconds())
	if accessTokenMaxAge < 0 {
		accessTokenMaxAge = 0
	}

	ctx.SetCookie("accessToken", res.RawAccessToken, accessTokenMaxAge, "/", "", true, true)

	refreshTokenMaxAge := int(res.RefreshTokenTTL.Seconds())
	if refreshTokenMaxAge < 0 {
		refreshTokenMaxAge = 0
	}

	ctx.SetCookie("refreshToken", res.RawRefreshToken, refreshTokenMaxAge, "/api/v1/auth/refresh", "", true, true)

	ctx.JSON(http.StatusOK, gin.H{
		"details": "Login successful",
		"user":    res.User,
	})
}

func (c *authController) RefreshToken(ctx *gin.Context) {
	var (
		res *dto.RefreshTokenResponseDto
		err error
	)

	refreshTokenString, err := ctx.Cookie("refreshToken")
	if err != nil {
		if err == http.ErrNoCookie {
			log.Printf("Refresh token cookie not found in request")
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Missing refresh token",
				"details": "Refresh token is required. Please log in.",
			})
			return
		}
		log.Printf("Error retrieving refresh token cookie: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": "Error processing refresh token.",
		})
		return
	}

	userAgent := ctx.Request.UserAgent()
	clientIP := ctx.ClientIP()

	res, err = c.authService.RefreshToken(refreshTokenString, userAgent, clientIP)
	if err != nil {
		log.Printf("Refresh error: %v", err)

		if errors.IsInvalidTokenError(err) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid Token",
				"message": "The provided refresh token is invalid or expired.",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Refresh failed",
			"message": "An internal error occured. Please try again later.",
		})
		return
	}

	accessTokenMaxAge := int(res.AccessTokenTTL.Seconds())
	if accessTokenMaxAge < 0 {
		accessTokenMaxAge = 0
	}

	ctx.SetCookie("accessToken", res.RawAccessToken, accessTokenMaxAge, "/", "", true, true)

	ctx.JSON(http.StatusOK, gin.H{
		"details": "Access token refreshed successfully",
	})
}
