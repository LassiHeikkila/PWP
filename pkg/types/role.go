package types

import (
	"strings"
)

type Role int8

const (
	RoleNone Role = 0
	RoleUser Role = 1 << iota
	RoleMaintainer
	RoleAdministrator
	RoleRoot
)

func HasRole(a, b Role) bool {
	return a&b != 0
}

func RoleFromString(s string) Role {
	switch strings.ToLower(s) {
	case "user":
		return RoleUser
	case "maintainer":
		return RoleMaintainer
	case "administrator":
		return RoleAdministrator
	case "root":
		return RoleRoot
	default:
		return RoleNone
	}
}

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleMaintainer:
		return "maintainer"
	case RoleAdministrator:
		return "administrator"
	case RoleRoot:
		return "root"
	}

	return "invalid"
}
