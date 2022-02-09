package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name           string `gorm:"unique,not null"`
	OrganizationID uint   `gorm:"not null"`
	Role           int8   `gorm:"not null"`
}
