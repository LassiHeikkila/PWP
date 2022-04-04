package types

import (
	"time"
)

type Schedule struct {
	SingleshotTasks []SingleshotTask `json:"singleshot"`
	PeriodicTasks   []PeriodicTask   `json:"periodically"`
	CronTasks       []CronTask       `json:"cron"`
}

type SingleshotTask struct {
	When time.Time `json:"when"` // time.Parse(time.RFC3339Nano, ...)
	What string    `json:"taskID"`
}

type PeriodicTask struct {
	Interval time.Duration `json:"every"` // anything supported by http://golang.org/pkg/time/#ParseDuration
	What     string        `json:"taskID"`
}

type CronTask struct {
	When string `json:"cron"` // anything supported by default by https://github.com/robfig/cron
	What string `json:"taskID"`
}
