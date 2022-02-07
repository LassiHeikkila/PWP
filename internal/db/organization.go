package db

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Users    []User    `gorm:"foreignKey:ID"`
	Machines []Machine `gorm:"foreignKey:ID"`
}
