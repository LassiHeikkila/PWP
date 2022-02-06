package db

import (
	"time"

	"gorm.io/gorm"
)

type MachineToken struct {
	gorm.Model
	MachineID  Machine
	Expiration time.Time
}

type UserToken struct {
	gorm.Model
	UserID     User
	Expiration time.Time
}
