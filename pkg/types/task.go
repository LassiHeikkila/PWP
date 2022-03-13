package types

type Task struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content"`
}
