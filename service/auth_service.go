package service

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/AyKrimino/ScribeHost/dto"
	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/AyKrimino/ScribeHost/errors"
	"github.com/AyKrimino/ScribeHost/helper"
	"github.com/AyKrimino/ScribeHost/repository"
)

type AuthService interface {
	Register(req dto.RegisterRequestDto) (*dto.RegisterResponseDto, error)
	Login(req dto.LoginRequestDto, userAgent, clientIP string) (*dto.LoginResponseDto, error)
	RefreshToken(req dto.RefreshTokenRequestDto) (*dto.RefreshTokenResponseDto, error)
}

type authService struct {
	authRepo         repository.AuthRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewAuthService(authRepo repository.AuthRepository, refreshTokenRepo repository.RefreshTokenRepository) AuthService {
	return &authService{
		authRepo:         authRepo,
		refreshTokenRepo: refreshTokenRepo,
	}
}

func (s *authService) Register(req dto.RegisterRequestDto) (*dto.RegisterResponseDto, error) {
	existingUser, err := s.authRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.NewObjectAlreadyExistsError("user", "email", req.Email)
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

func (s *authService) Login(req dto.LoginRequestDto, userAgent, clientIP string) (*dto.LoginResponseDto, error) {
	existingUser, err := s.authRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser == nil {
		return nil, errors.NewObjectNotFoundError("user", "email", req.Email)
	}

	isPasswordValid := helper.ComparePassword(existingUser.PasswordHash, []byte(req.Password))
	if !isPasswordValid {
		return nil, errors.NewInvalidCredentialsError("password")
	}

	if !existingUser.IsActive {
		return nil, errors.NewObjectNotActiveError("user")
	}

	// TODO: check if the email is verified

	accessToken, err := helper.CreateToken(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	refreshTokenString := rand.Text()
	refreshTokenHashed, err := helper.HashPassword(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	refreshToken := entity.RefreshToken{
		TokenHash: refreshTokenHashed,
		UserID:    existingUser.ID,
		Expiry:    time.Now().UTC().Add(time.Hour * 24 * 15),
		IssuedAt:  time.Now().UTC(),
		UserAgent: userAgent,
		IpAddress: clientIP,
	}

	err = s.refreshTokenRepo.Create(&refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to save refresh token: %w", err)
	}

	now := time.Now().UTC()
	existingUser.LastLogin = &now

	err = s.authRepo.Update(*existingUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user last login: %w", err)
	}

	res := dto.LoginResponseDto{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenString,
		TokenType:    "Bearer",
		User:         dto.FromEntityToRegisterResponseDto(existingUser),
		LoggedInAt:   now,
	}

	return &res, nil
}

func (s *authService) RefreshToken(req dto.RefreshTokenRequestDto) (*dto.RefreshTokenResponseDto, error) {
	return nil, nil
}
