package service

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/dto"
	"github.com/AyKrimino/ScribeHost/errors"
	"github.com/AyKrimino/ScribeHost/helper"
	"github.com/AyKrimino/ScribeHost/repository"
)

type AuthService interface {
	Register(req dto.RegisterRequestDto) (*dto.RegisterResponseDto, error)
	Login(req dto.LoginRequestDto) (*dto.LoginResponseDto, error)
	RefreshToken(req dto.RefreshTokenRequestDto) (*dto.RefreshTokenResponseDto, error)
}

type authService struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &authService{
		authRepo: authRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequestDto) (*dto.RegisterResponseDto, error) {
	existingUser, err := s.authRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser != nil {
		return nil, errors.NewObjectAlreadyExistsError("user", req.Email)
	}

	user := req.ToEntity()

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = hashed

	createdUser, err := s.authRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user to database: %w", err)
	}

	response := dto.FromEntityToRegisterResponseDto(createdUser)
	return &response, nil
}

func (s *authService) Login(req dto.LoginRequestDto) (*dto.LoginResponseDto, error) {
	return nil, nil
}

func (s *authService) RefreshToken(req dto.RefreshTokenRequestDto) (*dto.RefreshTokenResponseDto, error) {
	return nil, nil
}
