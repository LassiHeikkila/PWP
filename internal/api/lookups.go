package api

import (
	"github.com/jackc/pgtype"

	"github.com/LassiHeikkila/taskey/internal/auth"
	"github.com/LassiHeikkila/taskey/internal/db"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func lookupUserByJWT(authController auth.Controller, token string) *types.User {
	var user, organization string
	var role int
	if !authController.ValidateUserToken(
		token,
		&user,
		&organization,
		&role,
	) {
		return nil
	}
	return &types.User{
		Name:         user,
		Organization: organization,
		Role:         types.Role(role),
	}
}

func lookupUserByToken(dbController db.Controller, token string) *types.User {
	u := pgtype.UUID{}
	if err := u.Set(token); err != nil {
		return nil
	}

	r, err := dbController.ReadUserToken(u)
	if err != nil {
		return nil
	}
	return &types.User{
		Name:         r.User.Name,
		Organization: r.User.Organization.Name,
		Role:         r.User.Role,
	}
}

func lookupMachineByToken(dbController db.Controller, token string) *types.Machine {
	u := pgtype.UUID{}
	if err := u.Set(token); err != nil {
		return nil
	}

	r, err := dbController.ReadMachineToken(u)
	if err != nil {
		return nil
	}
	return &types.Machine{
		Name:        r.Machine.Name,
		Description: r.Machine.Description,
		OS:          r.Machine.OS,
		Arch:        r.Machine.Arch,
	}
}
