package db

import (
	"gorm.io/gorm"
)

type LoginInfo struct {
	gorm.Model
	Username string `gorm:"unique,not null"`
	Password string `gorm:"not null"`
	UserID   uint
	User     User
}
