package db

import (
	"gorm.io/gorm"
)

type LoginInfo struct {
	gorm.Model
	Username string
	Password string
	UserID   int
	User     User
}
