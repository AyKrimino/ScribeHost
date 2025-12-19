package repository

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/internal/entity"
	"github.com/jinzhu/gorm"
)

type RefreshTokenRepository interface {
	Create(token *entity.RefreshToken) error
	FindByTokenHash(tokenHash string) (*entity.RefreshToken, error)
	Revoke(tokenHash string) error
	RevokeAllForUser(userId uint) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{
		db: DB,
	}
}

func (r *refreshTokenRepository) Create(token *entity.RefreshToken) error {
	token.IssuedAt = token.IssuedAt.UTC()
	token.Expiry = token.Expiry.UTC()

	if err := r.db.Create(token).Error; err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepository) FindByTokenHash(tokenHash string) (*entity.RefreshToken, error) {
	var token entity.RefreshToken

	err := r.db.Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find refresh token by hash %s: %w", tokenHash, err)
	}
	return &token, nil
}

func (r *refreshTokenRepository) Revoke(tokenHash string) error {
	err := r.db.Model(&entity.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Updates(map[string]interface{}{
			"is_revoked": true,
		}).
		Error
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token with hash %s: %w", tokenHash, err)
	}
	return nil
}

func (r *refreshTokenRepository) RevokeAllForUser(userId uint) error {
	err := r.db.Model(&entity.RefreshToken{}).
		Where("user_id = ? AND is_revoked = ?", userId, false).
		Updates(map[string]interface{}{
			"is_revoked": true,
		}).
		Error
	if err != nil {
		return fmt.Errorf("failed to revoke refresh tokens for user %d: %w", userId, err)
	}
	return nil
}
