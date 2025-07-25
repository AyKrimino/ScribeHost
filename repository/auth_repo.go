package repository

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/jinzhu/gorm"
)

type AuthRepository interface {
	CreateUser(user entity.User) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
}

type authRepositoy struct {
	db *gorm.DB
}

func NewAuthRepository() AuthRepository {
	return &authRepositoy{
		db: DB,
	}
}

func (r *authRepositoy) CreateUser(user entity.User) (*entity.User, error) {
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

func (r *authRepositoy) FindUserByEmail(email string) (*entity.User, error) {
	var (
		user entity.User
		err  error
	)

	err = r.db.Where("email = ?", email).Preload("Profile").First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email %s: %w", email, err)
	}

	return &user, nil
}
