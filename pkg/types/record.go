package types

import (
	"time"
)

type Record struct {
	MachineName string
	TaskName    string
	ExecutedAt  time.Time
	Status      int
	Output      string
}
