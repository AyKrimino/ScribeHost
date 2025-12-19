package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/AyKrimino/ScribeHost/internal/entity"
	em "github.com/AyKrimino/ScribeHost/internal/infrastructure/email"
	infrajwt "github.com/AyKrimino/ScribeHost/internal/infrastructure/jwt"
)

type AuthService interface {
	Register(req RegisterRequestDto) (*RegisterResponseDto, error)
	Login(req LoginRequestDto, userAgent, clientIP string) (*LoginResponseDto, error)
	RefreshToken(refreshTokenString, userAgent, clientIP string) (*RefreshTokenResponseDto, error)
	Logout(userID uint) (*LogoutResponseDto, error)
	sendOTP(email string) error
	ResendOTP(email string) (*ResendOTPResponseDto, error)
	VerifyOTP(email, otp string) (*VerifyOTPResponseDto, error)
	sendResetPassword(email string) error
	ForgotPassword(email string) (*ForgotPasswordResponseDto, error)
	ResetPassword(email, token, newPassword string) (*ResetPasswordResponseDto, error)
}

type authService struct {
	authRepo               AuthRepository
	refreshTokenRepo       RefreshTokenRepository
	otpRedisRepo           OtpRedisRepo
	passwordResetRedisRepo PasswordResetRedisRepo
}

func NewAuthService(
	authRepo AuthRepository,
	refreshTokenRepo RefreshTokenRepository,
	otpRedisRepo OtpRedisRepo,
	passwordResetRedisRepo PasswordResetRedisRepo,
) AuthService {
	return &authService{
		authRepo:               authRepo,
		refreshTokenRepo:       refreshTokenRepo,
		otpRedisRepo:           otpRedisRepo,
		passwordResetRedisRepo: passwordResetRedisRepo,
	}
}

func (s *authService) Register(req RegisterRequestDto) (*RegisterResponseDto, error) {
	existingUser, err := s.authRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser != nil {
		return nil, NewObjectAlreadyExistsError("user", "email", req.Email)
	}

	user := req.ToEntity()

	hashed, err := HashPassword(req.Password)
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

	response := FromEntityToRegisterResponseDto(createdUser)
	return &response, nil
}

func (s *authService) Login(req LoginRequestDto, userAgent, clientIP string) (*LoginResponseDto, error) {
	existingUser, err := s.authRepo.FindUserByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if existingUser == nil {
		return nil, NewObjectNotFoundError("user", "email", req.Email)
	}

	isPasswordValid := ComparePassword(existingUser.PasswordHash, []byte(req.Password))
	if !isPasswordValid {
		return nil, NewInvalidCredentialsError("password")
	}

	if !existingUser.IsActive {
		return nil, NewObjectNotActiveError("user")
	}

	if !existingUser.EmailVerified {
		return nil, fmt.Errorf("email Verification required")
	}

	accessTokenString, err := infrajwt.CreateToken(existingUser.ID, existingUser.Email, existingUser.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	refreshTokenString := rand.Text()
	refreshTokenHashed, err := infrajwt.HashToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to hash refresh token: %w", err)
	}

	refreshToken := entity.RefreshToken{
		TokenHash: refreshTokenHashed,
		UserID:    existingUser.ID,
		Expiry:    time.Now().UTC().Add(time.Hour * 24 * 15), // 15 days
		IssuedAt:  time.Now().UTC(),
		UserAgent: userAgent,
		IPAddress: clientIP,
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

	parsedAccessToken, err := infrajwt.ValidateAndExtractClaims(accessTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated access token for expiry: %w", err)
	}
	_, accessTokenTTL, err := infrajwt.GetTokenExpiry(parsedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token expiry: %w", err)
	}

	refreshTokenTTL := 15 * 24 * time.Hour

	res := LoginResponseDto{
		RawAccessToken:  accessTokenString,
		RawRefreshToken: refreshTokenString,
		TokenType:       "Bearer",
		AccessTokenTTL:  accessTokenTTL,
		RefreshTokenTTL: refreshTokenTTL,
		User:            FromEntityToRegisterResponseDto(existingUser),
		LoggedInAt:      now,
	}

	return &res, nil
}

func (s *authService) RefreshToken(refreshTokenString, userAgent, clientIP string) (*RefreshTokenResponseDto, error) {
	hashed, err := infrajwt.HashToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to hash the refresh token: %w", err)
	}

	refreshToken, err := s.refreshTokenRepo.FindByTokenHash(hashed)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup refresh token: %w", err)
	}
	if refreshToken == nil {
		return nil, NewInvalidTokenError("refreshToken", "refresh token not found")
	}

	if !refreshToken.IsValid() {
		return nil, NewInvalidTokenError("refreshToken", "refresh token expired")
	}

	user, err := s.authRepo.FindUserById(refreshToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user associated with refresh token: %w", err)
	}
	if user == nil {
		return nil, NewInvalidTokenError("refreshToken", "associated user not found")
	}

	accessTokenString, err := infrajwt.CreateToken(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create token: %w", err)
	}

	parsedAccessToken, err := infrajwt.ValidateAndExtractClaims(accessTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated access token for expiry: %w", err)
	}
	_, accessTokenTTL, err := infrajwt.GetTokenExpiry(parsedAccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token expiry: %w", err)
	}

	res := RefreshTokenResponseDto{
		RawAccessToken: accessTokenString,
		TokenType:      "Bearer",
		AccessTokenTTL: accessTokenTTL,
	}

	return &res, nil
}

func (s *authService) Logout(userID uint) (*LogoutResponseDto, error) {
	err := s.refreshTokenRepo.RevokeAllForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to logout: unable to revoke session tokens")
	}

	return &LogoutResponseDto{
		Msg: "Successfully logged out",
	}, nil
}

func (s *authService) sendOTP(email string) error {
	otp := GenerateOTP()

	err := s.otpRedisRepo.StoreOTP(email, otp)
	if err != nil {
		return fmt.Errorf("failed to store OTP: %w", err)
	}

	err = em.SendOTPByEmail(email, otp)
	if err != nil {
		return fmt.Errorf("failed to send otp to email %s: %w", email, err)
	}

	return nil
}

func (s *authService) VerifyOTP(email, otp string) (*VerifyOTPResponseDto, error) {
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

	return &VerifyOTPResponseDto{
		Msg: "OTP verified successfully",
	}, nil
}

func (s *authService) ResendOTP(email string) (*ResendOTPResponseDto, error) {
	err := s.sendOTP(email)
	if err != nil {
		return &ResendOTPResponseDto{
			Msg: fmt.Sprintf("failed to resend OTP to email %s: %v", email, err),
		}, fmt.Errorf("failed to resend OTP to email %s: %w", email, err)
	}

	return &ResendOTPResponseDto{
		Msg: fmt.Sprintf("OTP successfully resent to email %s", email),
	}, nil
}

func (s *authService) sendResetPassword(email string) error {
	token, err := GeneratePasswordResetToken()
	if err != nil {
		return fmt.Errorf("failed to generate password reset token: %w", err)
	}

	resetLink := fmt.Sprintf("/reset-password?token=%s", token)

	err = s.passwordResetRedisRepo.StorePasswordResetToken(email, token)
	if err != nil {
		return fmt.Errorf("failed to store token: %w", err)
	}

	err = em.SendPasswordResetEmail(email, resetLink)
	if err != nil {
		return fmt.Errorf("failed to send token to email %s: %w", email, err)
	}

	return nil
}

func (s *authService) ForgotPassword(email string) (*ForgotPasswordResponseDto, error) {
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

	return &ForgotPasswordResponseDto{
		Msg: "reset password email was sent successfully",
	}, nil
}

func (s *authService) ResetPassword(email, token, newPassword string) (*ResetPasswordResponseDto, error) {
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

	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash the new password: %w", err)
	}

	user.PasswordHash = hashedPassword

	err = s.authRepo.Update(*user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user with the new password: %w", err)
	}

	return &ResetPasswordResponseDto{
		Msg: "new password is reset successfully",
	}, nil
}
