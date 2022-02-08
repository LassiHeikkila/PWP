package db

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Schedule struct {
	gorm.Model
	Content   pgtype.JSON
	MachineID uint
	Machine   Machine
}
