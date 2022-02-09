package db

import (
	"time"

	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type MachineToken struct {
	gorm.Model
	Value      pgtype.UUID `gorm:"type:uuid"`
	Expiration time.Time
	MachineID  uint
	Machine    Machine
}

type UserToken struct {
	gorm.Model
	Value      pgtype.UUID `gorm:"type:uuid"`
	Expiration time.Time
	UserID     uint
	User       User
}
