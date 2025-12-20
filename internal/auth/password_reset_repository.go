package auth

import (
	"context"
	"fmt"
	"time"

	infrajwt "github.com/AyKrimino/ScribeHost/internal/infra/jwt"
	"github.com/redis/go-redis/v9"
)

type PasswordResetRedisRepo interface {
	StorePasswordResetToken(email string, token string) error
	VerifyPasswordResetToken(email string, token string) error
	deletePasswordResetToken(key string) error
}

type passwordResetRedisRepo struct {
	client *redis.Client
	ctx    context.Context
}

func NewPasswordResetRedisRepo(client *redis.Client, ctx context.Context) PasswordResetRedisRepo {
	return &passwordResetRedisRepo{
		client: client,
		ctx:    ctx,
	}
}

func (r *passwordResetRedisRepo) StorePasswordResetToken(email string, token string) error {
	hashedToken, err := infrajwt.HashToken(token)
	if err != nil {
		return fmt.Errorf("failed to hash token: %w", err)
	}

	key := fmt.Sprintf("password_reset:email:%s", email)

	err = r.client.Set(r.ctx, key, hashedToken, time.Hour).Err()

	return err
}

func (r *passwordResetRedisRepo) VerifyPasswordResetToken(email string, token string) error {
	hashedToken, err := infrajwt.HashToken(token)
	if err != nil {
		return fmt.Errorf("failed to hash token: %w", err)
	}

	key := fmt.Sprintf("password_reset:email:%s", email)
	storedToken, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("token with this key %s not found: %w", key, err)
	}

	if hashedToken != storedToken {
		return fmt.Errorf("tokens didn't match")
	}

	err = r.deletePasswordResetToken(key)
	if err != nil {
		return fmt.Errorf("failed to delete password reset token: %w", err)
	}

	return nil
}

func (r *passwordResetRedisRepo) deletePasswordResetToken(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	return err
}
