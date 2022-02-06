package db

import (
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	Users    []User
	Machines []Machine
}
