package api

// sign up has to support:
// - creation of an organization with initial user
// - creation of new user accounts by org admin
// - new users can set their own password when logging in the first time

type CreateOrganization struct {
	OrganizationName string `json:"organizationName"`
	AdminName        string `json:"adminUsername"`
	AdminPassword    string `json:"adminPassword"`
}

type CreateOrganizationAccount struct {
	OrganizationName string   `json:"organizationName"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	Roles            []string `json:"roles"`
}
