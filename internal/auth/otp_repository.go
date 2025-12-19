package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/AyKrimino/ScribeHost/helper"
	"github.com/redis/go-redis/v9"
)

type OtpRedisRepo interface {
	StoreOTP(email, otp string) error
	VerifyOTP(email, otp string) error
	deleteOTP(email string) error
}

type otpRedisRepo struct {
	client *redis.Client
	ctx    context.Context
}

func NewOtpRedisRepo(client *redis.Client, ctx context.Context) OtpRedisRepo {
	return &otpRedisRepo{
		client: client,
		ctx:    ctx,
	}
}

func (r *otpRedisRepo) StoreOTP(email, otp string) error {
	hashedOTP := helper.HashOTP(otp)
	err := r.client.Set(r.ctx, email, hashedOTP, 5*time.Minute).Err()

	return err
}

func (r *otpRedisRepo) VerifyOTP(email, otp string) error {
	hashedOTP := helper.HashOTP(otp)
	storedOTP, err := r.client.Get(r.ctx, email).Result()
	if err != nil {
		return NewInvalidOTPError("email key not found.")
	}

	if storedOTP != hashedOTP {
		return NewInvalidOTPError("otp didn't match.")
	}

	err = r.deleteOTP(email)
	if err != nil {
		return fmt.Errorf("failed to delete otp: %w", err)
	}

	return nil
}

func (r *otpRedisRepo) deleteOTP(email string) error {
	err := r.client.Del(r.ctx, email).Err()
	return err
}
