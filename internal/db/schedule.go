package db

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Machine Machine
	Content pgtype.JSON
}
