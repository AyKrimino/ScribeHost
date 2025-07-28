package dto

import (
	"time"

	"github.com/AyKrimino/ScribeHost/entity"
)

type RegisterRequestDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type RegisterResponseDto struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequestDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponseDto struct {
	RawAccessToken  string              `json:"-"`
	RawRefreshToken string              `json:"-"`
	TokenType       string              `json:"token_type"` // Bearer token
	User            RegisterResponseDto `json:"user"`
	LoggedInAt      time.Time           `json:"logged_in_at"`
	AccessTokenTTL  time.Duration       `json:"access_token_time_to_live"`
	RefreshTokenTTL time.Duration       `json:"refresh_token_time_to_live"`
}

type RefreshTokenResponseDto struct {
	RawAccessToken string        `json:"-"`
	AccessTokenTTL time.Duration `json:"access_token_time_to_live"`
	TokenType      string        `json:"token_type"`
}

func (req *RegisterRequestDto) ToEntity() entity.User {
	return entity.User{
		Email: req.Email,
	}
}

func FromEntityToRegisterResponseDto(user *entity.User) RegisterResponseDto {
	return RegisterResponseDto{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
