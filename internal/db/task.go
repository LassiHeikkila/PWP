package db

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name        string
	Description string
	Content     pgtype.JSON
}
