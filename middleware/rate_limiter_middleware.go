package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AyKrimino/ScribeHost/dto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimitState struct {
	Tokens     int
	LastRefill int64
}

type Bucket struct {
	BucketSize   int
	RefillRate   int // requests per hour
	rateLimitKey string
}

var RateLimits = map[string]Bucket{
	"register":   {BucketSize: 1, RefillRate: 1, rateLimitKey: "ip"},
	"login":      {BucketSize: 5, RefillRate: 5, rateLimitKey: "ip"},
	"verify-otp": {BucketSize: 3, RefillRate: 3, rateLimitKey: "email"},
	"resend-otp": {BucketSize: 3, RefillRate: 3, rateLimitKey: "ip"},
}

var redisClient *redis.Client

func SetRedisClient(client *redis.Client) {
	redisClient = client
}

func RateLimiter(endpoint string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if redisClient == nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiting not configured",
			})
			return
		}

		var identifier string
		switch endpoint {
		case "verify-otp":
			var req dto.VerifyOTPRequestDto
			if err := ctx.ShouldBindJSON(&req); err != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid input",
					"details": "The request data is invalid. Please check your input and try again.",
				})
				return
			}
			identifier = req.Email
		default:
			identifier = ctx.ClientIP()
		}

		rateLimit, ok := RateLimits[endpoint]
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "Invalid endpoint",
			})
			return
		}

		key := fmt.Sprintf("rate_limit:%s:%s:%s", endpoint, rateLimit.rateLimitKey, identifier)

		c := context.Background()
		val, err := redisClient.Get(c, key).Result()
		if err != nil && err != redis.Nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check rate limit",
			})
			return
		}

		var state RateLimitState
		if err == redis.Nil {
			state = RateLimitState{
				Tokens:     rateLimit.BucketSize,
				LastRefill: time.Now().Unix(),
			}
		} else {
			if err := json.Unmarshal([]byte(val), &state); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Failed to parse rate limit state",
				})
				return
			}
		}

		now := time.Now().Unix()
		timeSinceLastRefill := now - state.LastRefill
		tokensToAdd := int(float64(timeSinceLastRefill) * float64(rateLimit.RefillRate) / 3600.0)
		newTokens := min(state.Tokens+tokensToAdd, rateLimit.BucketSize)

		if newTokens <= 0 {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			return
		}

		state.Tokens = newTokens - 1
		state.LastRefill = now

		stateBytes, err := json.Marshal(state)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update rate limit state",
			})
			return
		}

		err = redisClient.Set(c, key, stateBytes, 2*time.Hour).Err()
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update rate limit state",
			})
			return
		}

		ctx.Next()
	}
}
