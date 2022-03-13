package types

type User struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
	Role         Role   `json:"role"`
}
