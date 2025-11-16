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
	VerifyOTP(email, otp string) (*dto.VerifyOTPResponseDto, error)
	sendResetPassword(email string) error
	ForgotPassword(email string) (*dto.ForgotPasswordResponseDto, error)
	ResetPassword(email, token, newPassword string) (*dto.ResetPasswordResponseDto, error)
}

type authService struct {
	authRepo               repository.AuthRepository
	refreshTokenRepo       repository.RefreshTokenRepository
	otpRedisRepo           repository.OtpRedisRepo
	passwordResetRedisRepo repository.PasswordResetRedisRepo
}

func NewAuthService(
	authRepo repository.AuthRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	otpRedisRepo repository.OtpRedisRepo,
	passwordResetRedisRepo repository.PasswordResetRedisRepo,
) AuthService {
	return &authService{
		authRepo:               authRepo,
		refreshTokenRepo:       refreshTokenRepo,
		otpRedisRepo:           otpRedisRepo,
		passwordResetRedisRepo: passwordResetRedisRepo,
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
	otp := helper.GenerateOTP()

	err := s.otpRedisRepo.StoreOTP(email, otp)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	err = helper.SendOTPByEmail(email, otp)
	if err != nil {
		return fmt.Errorf("failed to send otp to email %s: %w", email, err)
	}

	return nil
}

func (s *authService) VerifyOTP(email, otp string) (*dto.VerifyOTPResponseDto, error) {
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user with email %s: %w", email, err)
	}

	err = s.otpRedisRepo.VerifyOTP(email, otp)
	if err != nil {
		return nil, err
	}

	user.EmailVerified = true
	err = s.authRepo.Update(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
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

func (s *authService) sendResetPassword(email string) error {
	token, err := helper.GeneratePasswordResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate password reset token: %w", err)
	}

	resetLink := fmt.Sprintf("/reset-password?token=%s", token)

	err = s.passwordResetRedisRepo.StorePasswordResetToken(email, token)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	err = helper.SendPasswordResetEmail(email, resetLink)
	if err != nil {
		return fmt.Errorf("failed to send token to email %s: %w", email, err)
	}

	return nil
}

func (s *authService) ForgotPassword(email string) (*dto.ForgotPasswordResponseDto, error) {
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user with email %s: %w", email, err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with email %s does not exist", email)
	}

	err = s.sendResetPassword(email)
	if err != nil {
		return nil, err
	}

	return &dto.ForgotPasswordResponseDto{
		Msg: "reset password email was sent successfully",
	}, nil
}

func (s *authService) ResetPassword(email, token, newPassword string) (*dto.ResetPasswordResponseDto, error) {
	err := s.passwordResetRedisRepo.VerifyPasswordResetToken(email, token)
	if err != nil {
		return nil, fmt.Errorf("failed to verify password reset token: %w", err)
	}

	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user with email %s: %w", email, err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with email %s does not exist", email)
	}

	hashedPassword, err := helper.HashPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash the new password: %w", err)
	}

	user.PasswordHash = hashedPassword

	err = s.authRepo.Update(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user with the new password: %w", err)
	}

	return &dto.ResetPasswordResponseDto{
		Msg: "new password is reset successfully",
	}, nil
}
