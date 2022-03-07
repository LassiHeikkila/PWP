package dbconverter

import (
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func ConvertUser(dbuser db.User) types.User {
	return types.User{
		Name:         dbuser.Name,
		Email:        dbuser.Email,
		Organization: dbuser.Organization.Name,
		Role:         dbuser.Role,
	}
}
