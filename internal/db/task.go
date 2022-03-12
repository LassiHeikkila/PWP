package db

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Name           string `gorm:"unique,not null"`
	Description    string
	Content        pgtype.JSON `gorm:"type:json"`
	Records        []Record    `gorm:"foreignKey:TaskID"`
	OrganizationID uint        `gorm:"not null"`
}

func StringToJSON(s string) pgtype.JSON {
	j := pgtype.JSON{}
	j.Set(s)
	return j
}
