package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/AyKrimino/ScribeHost/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header missing",
			})
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := helper.ValidateAndExtractClaims(tokenString)
		if err != nil {
			log.Printf("JWT Middleware: Token validation failed: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("JWT Middleware: Failed to assert claims")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token claims",
			})
			return
		}

		userID, err := helper.GetUserIDFromClaims(claims)
		if err != nil {
			log.Printf("JWT Middleware: Failed to extract user ID: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid user ID in token",
			})
			return
		}

		role, err := helper.GetRoleFromClaims(claims)
		if err != nil {
			log.Printf("JWT Middleware: Failed to extract role: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid role in token",
			})
			return
		}

		email, err := helper.GetEmailFromClaims(claims)
		if err != nil {
			log.Printf("JWT Middleware: Failed to extract email: %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email in token",
			})
			return
		}

		ctx.Set("token", token)
		ctx.Set("userID", userID)
		ctx.Set("userRole", role)
		ctx.Set("userEmail", email)

		ctx.Next()
	}
}
