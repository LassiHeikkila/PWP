package types

type Role int8

const (
	RoleUser Role = 1 << iota
	RoleMaintainer
	RoleAdministrator
	RoleRoot
)

func HasRole(a, b Role) bool {
	return a&b != 0
}
