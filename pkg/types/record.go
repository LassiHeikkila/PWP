package types

import (
	"time"
)

type Record struct {
	MachineName string    `json:"machineName"`
	TaskName    string    `json:"taskName"`
	ExecutedAt  time.Time `json:"executedAt"`
	Status      int       `json:"status"`
	Output      string    `json:"output"`
}
