package db

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type MachineToken struct {
	gorm.Model
	Value      pgtype.UUID
	MachineID  Machine `gorm:"foreignKey:ID"`
	Expiration time.Time
}

type UserToken struct {
	gorm.Model
	Value      pgtype.UUID
	UserID     User `gorm:"foreignKey:ID"`
	Expiration time.Time
}
