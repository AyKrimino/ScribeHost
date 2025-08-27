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
	RefreshToken(refreshTokenString, userAgent, clientIP string) (*dto.RefreshTokenResponseDto, error)
	Logout(userID uint) (*dto.LogoutResponseDto, error)
	sendOTP(email string) error
	ResendOTP(email string) (*dto.ResendOTPResponseDto, error)
	StoreOTP(user *entity.User, otp string) error
	VerifyOTP(email, otp string) (*dto.VerifyOTPResponseDto, error)
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

	err = s.sendOTP(createdUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to send otp: %w", err)
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

	if !existingUser.EmailVerified {
		return nil, fmt.Errorf("Email Verification required")
	}

	accessTokenString, err := helper.CreateToken(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	refreshTokenString := rand.Text()
	refreshTokenHashed, err := helper.HashToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	refreshToken := entity.RefreshToken{
		TokenHash: refreshTokenHashed,
		UserID:    existingUser.ID,
		Expiry:    time.Now().UTC().Add(time.Hour * 24 * 15), // 15 days
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

	parsedAccessToken, err := helper.ValidateAndExtractClaims(accessTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated access token for expiry: %w", err)
	}
	_, accessTokenTTL, err := helper.GetTokenExpiry(parsedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token expiry: %w", err)
	}

	refreshTokenTTL := 15 * 24 * time.Hour

	res := dto.LoginResponseDto{
		RawAccessToken:  accessTokenString,
		RawRefreshToken: refreshTokenString,
		TokenType:       "Bearer",
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		User:            dto.FromEntityToRegisterResponseDto(existingUser),
		LoggedInAt:      now,
	}

	return &res, nil
}

func (s *authService) RefreshToken(refreshTokenString, userAgent, clientIP string) (*dto.RefreshTokenResponseDto, error) {
	hashed, err := helper.HashToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to hash the refresh token: %w", err)
	}

	refreshToken, err := s.refreshTokenRepo.FindByTokenHash(hashed)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup refresh token: %w", err)
	}
	if refreshToken == nil {
		return nil, errors.NewInvalidTokenError("refreshToken", "refresh token not found")
	}

	if !refreshToken.IsValid() {
		return nil, errors.NewInvalidTokenError("refreshToken", "refresh token expired")
	}

	user, err := s.authRepo.FindUserById(refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user associated with refresh token: %w", err)
	}
	if user == nil {
		return nil, errors.NewInvalidTokenError("refreshToken", "associated user not found")
	}

	accessTokenString, err := helper.CreateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	parsedAccessToken, err := helper.ValidateAndExtractClaims(accessTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated access token for expiry: %w", err)
	}
	_, accessTokenTTL, err := helper.GetTokenExpiry(parsedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token expiry: %w", err)
	}

	res := dto.RefreshTokenResponseDto{
		RawAccessToken: accessTokenString,
		TokenType:      "Bearer",
		AccessTokenTTL: accessTokenTTL,
	}

	return &res, nil
}

func (s *authService) Logout(userID uint) (*dto.LogoutResponseDto, error) {
	err := s.refreshTokenRepo.RevokeAllForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to logout: unable to revoke session tokens")
	}

	return &dto.LogoutResponseDto{
		Msg: "Successfully logged out",
	}, nil
}

func (s *authService) sendOTP(email string) error {
	user, err := s.authRepo.FindUserByEmail(email)
	if user == nil {
		return fmt.Errorf("failed to find user with email %s: %w", email, err)
	}

	otp := helper.GenerateOTP()
	err = s.StoreOTP(user, otp)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	err = helper.SendOTPByEmail(user.Email, otp)
	if err != nil {
		return fmt.Errorf("failed to send otp to email %s: %w", user.Email, err)
	}

	return nil
}

func (s *authService) StoreOTP(user *entity.User, otp string) error {
	hashedOTP := helper.HashOTP(otp)
	user.OTP = hashedOTP

	now := time.Now()
	expiry := now.Add(5 * time.Minute)
	user.OTPExpiry = &expiry

	err := s.authRepo.Update(*user)
	return err
}

func (s *authService) VerifyOTP(email, otp string) (*dto.VerifyOTPResponseDto, error) {
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	hashedOTP := helper.HashOTP(otp)
	if hashedOTP != user.OTP {
		return nil, errors.NewInvalidOTPError("OTPs didn't match")
	}

	if user.OTPExpiry.Before(time.Now()) {
		return nil, errors.NewInvalidOTPError("OTP has been expired")
	}

	user.EmailVerified = true
	user.OTP = ""
	user.OTPExpiry = nil

	err = s.authRepo.Update(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user with OTP attributes updates")
	}

	return &dto.VerifyOTPResponseDto{
		Msg: "OTP verified successfully",
	}, nil
}

func (s *authService) ResendOTP(email string) (*dto.ResendOTPResponseDto, error) {
	err := s.sendOTP(email)
	if err != nil {
		return &dto.ResendOTPResponseDto{
			Msg: fmt.Sprintf("failed to resend OTP to email %s: %v", email, err),
		}, fmt.Errorf("failed to resend OTP to email %s: %w", email, err)
	}

	return &dto.ResendOTPResponseDto{
		Msg: fmt.Sprintf("OTP successfully resent to email %s", email),
	}, nil
}
