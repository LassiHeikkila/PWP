package db

import (
	"time"

	"gorm.io/gorm"
)

type Record struct {
	gorm.Model
	MachineID  uint `gorm:"not null"`
	Machine    Machine
	TaskID     uint `gorm:"not null"`
	Task       Task
	ExecutedAt time.Time
	Status     int
	Output     string
}
