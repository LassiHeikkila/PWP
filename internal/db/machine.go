package db

import (
	"gorm.io/gorm"
)

type Machine struct {
	gorm.Model
	Name           string `gorm:"unique,not null"`
	Description    string
	OS             string
	Arch           string
	OrganizationID uint
	ScheduleID     uint
	Records        []Record `gorm:"foreignKey:ID"`
}
