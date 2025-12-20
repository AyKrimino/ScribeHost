package user

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/internal/infra/crypto"
)

type UserService interface {
	CreateUser(userReq CreateUserRequestDto) (*UserResponseDto, error)
}

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(userReq CreateUserRequestDto) (*UserResponseDto, error) {
	user, err := userReq.ToEntity()
	if err != nil {
		return nil, fmt.Errorf("conversion error: %w", err)
	}

	hashed, err := crypto.HashPassword(userReq.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashed

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user to database: %w", err)
	}

	response := FromEntity(createdUser)
	return &response, nil
}
