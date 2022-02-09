package db

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Name     string    `gorm:"unique,not null"`
	Users    []User    `gorm:"foreignKey:OrganizationID"`
	Machines []Machine `gorm:"foreignKey:OrganizationID"`
}
