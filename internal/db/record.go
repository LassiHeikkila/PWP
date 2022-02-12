package db

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	MachineID uint
	TaskID    uint
	Timestamp time.Time
	Status    int
	Output    string
}
