package db

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name        string `gorm:"unique,not null"`
	Description string
	Content     pgtype.JSON
	Records     []Record `gorm:"foreignKey:TaskID"`
}
