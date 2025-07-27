package entity

import (
	"time"
)

type RefreshToken struct {
	TokenHash string    `gorm:"type:varchar(64);primary_key"`
	UserID    uint      `gorm:"index;not null;refrences:users(id)"`
	Expiry    time.Time `gorm:"index;not null"`
	IssuedAt  time.Time `gorm:"not null"`
	UserAgent string    `gorm:"type:text"`
	IpAddress string    `gorm:"type:varchar(45)"`
}
