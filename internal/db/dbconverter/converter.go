package dbconverter

import (
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func ConvertOrganization(dborg db.Organization) types.Organization {
	return types.Organization{
		Name: dborg.Name,
	}
}

func ConvertUser(dbuser db.User) types.User {
	return types.User{
		Name:         dbuser.Name,
		Email:        dbuser.Email,
		Organization: dbuser.Organization.Name,
		Role:         dbuser.Role,
	}
}

func ConvertMachine(dbmachine db.Machine) types.Machine {
	return types.Machine{
		Name:        dbmachine.Name,
		Description: dbmachine.Description,
		OS:          dbmachine.OS,
		Arch:        dbmachine.Arch,
	}
}
