package repository

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/jinzhu/gorm"
)

type UserRepository interface {
	Create(user entity.User) (*entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: DB,
	}
}

func (r *userRepository) Create(user entity.User) (*entity.User, error) {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	var createdUser entity.User
	if err := r.db.Preload("Profile").First(&createdUser, user.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload user: %w", err)
	}

	return &createdUser, nil
}
