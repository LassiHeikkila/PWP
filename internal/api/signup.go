package api

// sign up has to support:
// - creation of an organization with initial user
// - creation of new user accounts by org admin
// - new users can set their own password when logging in the first time

type SignUpRequest struct {
	OrganizationName string `json:"orgName"`
	Username         string `json:"username"`
	Email            string `json:"email"`
	Password         string `json:"password"`
}

func validateSignUpRequest(r *SignUpRequest) bool {
	return r.OrganizationName != "" && r.Username != "" && r.Email != "" && r.Password != ""
}
