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
	AccessToken  string              `json:"access_token"`
	RefreshToken string              `json:"refresh_token"`
	TokenType    string              `json:"token_type"` // Bearer token
	ExpiresIn    int                 `json:"expires_in"` // in seconds
	User         RegisterResponseDto `json:"user"`
	LoggedInAt   time.Time           `json:"logged_in_at"`
}

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshTokenResponseDto struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
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
