package types

type Task struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Content     interface{} `json:"content"`
}
