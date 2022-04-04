package types

type Task struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Content     interface{} `json:"content"` // CmdTask | ScriptTask
}

const (
	TaskTypeCmd    = "cmd"
	TaskTypeScript = "script"
)

type TaskProperties struct {
	Type           string `json:"type"`
	CombinedOutput bool   `json:"combinedOutput"`
}

type CmdTask struct {
	TaskProperties
	Program string   `json:"program"`
	Args    []string `json:"args"`
}

type ScriptTask struct {
	TaskProperties
	Interpreter string `json:"interpreter"` // sh, bash, zsh, python, etc.
	Script      string `json:"script"`
}
