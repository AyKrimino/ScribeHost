package repository

import (
	"fmt"

	"github.com/AyKrimino/ScribeHost/entity"
	"github.com/jinzhu/gorm"
)

type AuthRepository interface {
	CreateUser(user entity.User) (*entity.User, error)
	FindUserByEmail(email string) (*entity.User, error)
	FindUserById(id uint) (*entity.User, error)
	Update(user entity.User) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository() AuthRepository {
	return &authRepository{
		db: DB,
	}
}

func (r *authRepository) CreateUser(user entity.User) (*entity.User, error) {
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

func (r *authRepository) FindUserByEmail(email string) (*entity.User, error) {
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

func (r *authRepository) Update(user entity.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *authRepository) FindUserById(id uint) (*entity.User, error) {
	var (
		user entity.User
		err  error
	)

	err = r.db.Where("id = ?", id).Preload("Profile").First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id %d: %w", id, err)
	}
	return &user, nil
}
