package entity

import (
	"time"

	"github.com/AyKrimino/ScribeHost/types"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email         string         `gorm:"type:varchar(100);unique_index;not null" json:"email"`
	PasswordHash  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role          types.RoleType `gorm:"type:varchar(50);default:'author';index" json:"role"`
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	IsActive      bool           `gorm:"default:true;index" json:"is_active"`

	OTPSecret string     `gorm:"type:varchar(100)" json:"-"` // For 2FA
	OTP       string     `gorm:"type:varchar(10)" json:"-"`  // Current OTP code
	OTPExpiry *time.Time `json:"otp_expiry"`                 // OTP expiration

	// Password Reset
	ResetToken       string     `gorm:"type:varchar(100)" json:"-"`
	ResetTokenExpiry *time.Time `json:"reset_token_expiry"`

	RefreshToken       string     `gorm:"type:text" json:"-"`
	RefreshTokenExpiry *time.Time `json:"-"`

	// Timestamps
	LastLogin *time.Time `json:"last_login,omitempty"`

	// Relations
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"profile"`
}

type Profile struct {
	gorm.Model
	UserID      uint             `gorm:"unique_index;not null" json:"user_id"`
	FirstName   string           `gorm:"type:varchar(100)" json:"first_name"`
	LastName    string           `gorm:"type:varchar(100)" json:"last_name"`
	AvatarURL   string           `gorm:"type:text" json:"avatar_url"`
	Phone       string           `gorm:"type:varchar(20);index" json:"phone"`
	Address     string           `gorm:"type:text" json:"address"`
	DateOfBirth *time.Time       `json:"date_of_birth"`
	Gender      types.GenderType `gorm:"type:varchar(20);default:'male'" json:"gender"`
}
