package helper

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func CreateRedisClientConn(ctx context.Context) (*redis.Client, error) {
	LoadEnv()

	var (
		redisHost     = os.Getenv("REDIS_HOST")
		redisPort     = os.Getenv("REDIS_PORT")
		redisPassword = os.Getenv("REDIS_PASSWORD")
		redisDB       = os.Getenv("REDIS_DB")
	)

	redisDBi, err := strconv.Atoi(redisDB)
	if err != nil {
		return nil, fmt.Errorf("Error converting redisDB value to integer: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       redisDBi,
	})

	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("Unable to connect to redis: %w", err)
	}

	return client, nil
}
