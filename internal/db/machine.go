package db

import (
	"gorm.io/gorm"
)

type Machine struct {
	gorm.Model
	Name           string
	Description    string
	OS             string
	Arch           string
	OrganizationID uint
	ScheduleID     uint
}
