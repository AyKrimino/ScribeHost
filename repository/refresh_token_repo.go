package repository

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/jinzhu/gorm"
)

type RefreshTokenRepository interface {
	Create(token *entity.RefreshToken) error
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
