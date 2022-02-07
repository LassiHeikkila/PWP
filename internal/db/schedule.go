package db

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Machine Machine `gorm:"foreignKey:ID"`
	Content pgtype.JSON
}
