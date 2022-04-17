package types

import (
	"time"
)

type Record struct {
	ID          uint      `json:"id"`
	MachineName string    `json:"machineName,omitempty"`
	TaskName    string    `json:"taskName"`
	ExecutedAt  time.Time `json:"executedAt"`
	Status      int       `json:"status"`
	Output      string    `json:"output"`
}
