package dto

import (
	"time"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/AyKrimino/ScribeHost/types"
)

type CreateUserRequestDto struct {
	Email    string                  `json:"email" binding:"required,email"`
	Password string                  `json:"password" binding:"required,min=8"`
	Profile  CreateProfileRequestDto `json:"profile" binding:"required"`
}

type CreateProfileRequestDto struct {
	FirstName   string           `json:"first_name" binding:"required"`
	LastName    string           `json:"last_name" binding:"required"`
	Phone       string           `json:"phone"`
	Address     string           `json:"address"`
	DateOfBirth *time.Time       `json:"date_of_birth"`
	Gender      types.GenderType `json:"gender" binding:"omitempty,oneof=male female"`
	AvatarURL   string           `json:"avatar_url" binding:"omitempty,url"`
}

func (req *CreateUserRequestDto) ToEntity() (entity.User, error) {
	user := entity.User{
		Email: req.Email,
		Profile: entity.Profile{
			FirstName:   req.Profile.FirstName,
			LastName:    req.Profile.LastName,
			Phone:       req.Profile.Phone,
			Address:     req.Profile.Address,
			DateOfBirth: req.Profile.DateOfBirth,
			Gender:      req.Profile.Gender,
			AvatarURL:   req.Profile.AvatarURL,
		},
	}
	return user, nil
}

type UserResponseDto struct {
	ID        uint               `json:"id"`
	Email     string             `json:"email"`
	Role      types.RoleType     `json:"role"`
	IsActive  bool               `json:"is_active"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	Profile   ProfileResponseDto `json:"profile"`
}

type ProfileResponseDto struct {
	ID          uint             `json:"id"`
	UserID      uint             `json:"user_id"`
	FirstName   string           `json:"first_name"`
	LastName    string           `json:"last_name"`
	AvatarURL   string           `json:"avatar_url,omitempty"`
	Phone       string           `json:"phone,omitempty"`
	Address     string           `json:"address,omitempty"`
	DateOfBirth *time.Time       `json:"date_of_birth,omitempty"`
	Gender      types.GenderType `json:"gender"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

func FromEntity(user *entity.User) UserResponseDto {
	return UserResponseDto{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Profile: ProfileResponseDto{
			ID:          user.Profile.ID,
			UserID:      user.Profile.UserID,
			FirstName:   user.Profile.FirstName,
			LastName:    user.Profile.LastName,
			AvatarURL:   user.Profile.AvatarURL,
			Phone:       user.Profile.Phone,
			Address:     user.Profile.Address,
			DateOfBirth: user.Profile.DateOfBirth,
			Gender:      user.Profile.Gender,
			CreatedAt:   user.Profile.CreatedAt,
			UpdatedAt:   user.Profile.UpdatedAt,
		},
	}
}
