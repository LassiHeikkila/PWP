package types

type Machine struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
}
