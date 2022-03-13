package db

import (
	"gorm.io/gorm"

	"github.com/LassiHeikkila/taskey/pkg/types"
)

type User struct {
	gorm.Model
	Name           string     `gorm:"unique,not null"`
	Email          string     `gorm:"unique,not null"`
	OrganizationID uint       `gorm:"not null"`
	Role           types.Role `gorm:"not null"`
}
