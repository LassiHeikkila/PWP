package main

import (
	"errors"
	"log"
	"time"

	"github.com/LassiHeikkila/taskey/pkg/json"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func mustParseTime(in string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, in)
	if err != nil {
		panic("error parsing time: " + in + " as RFC3339Nano")
	}
	return t
}
func mustParseDuration(in string) time.Duration {
	d, err := time.ParseDuration(in)
	if err != nil {
		panic("error parsing duration: " + in)
	}
	return d
}

var (
	dummySchedule = &types.Schedule{
		SingleshotTasks: []types.SingleshotTask{
			types.SingleshotTask{
				What: "task123",
				When: mustParseTime("2022-04-30T10:04:00+03:00"),
			},
		},
		PeriodicTasks: []types.PeriodicTask{
			types.PeriodicTask{
				What:     "task456",
				Interval: json.Duration{mustParseDuration("3m")},
			},
		},
		CronTasks: []types.CronTask{
			types.CronTask{
				What: "task789",
				When: "0 */5 * * * *",
			},
		},
	}

	dummyTasks = map[string]*types.Task{
		"task123": &types.Task{
			Name:        "task123",
			Description: "example task 123",
			Content: &types.ScriptTask{
				TaskProperties: types.TaskProperties{
					Type:           types.TaskTypeScript,
					CombinedOutput: false,
				},
				Interpreter: "bash",
				Script: `echo "running script"
                exit 1`,
			},
		},
		"task456": &types.Task{
			Name:        "task456",
			Description: "example task 456",
			Content: &types.CmdTask{
				TaskProperties: types.TaskProperties{
					Type:           types.TaskTypeCmd,
					CombinedOutput: false,
				},
				Program: "/usr/bin/curl",
				Args:    []string{"https://taskey-service.herokuapp.com/api/v1/health/"},
			},
		},
		"task789": &types.Task{
			Name:        "task789",
			Description: "example task 789",
			Content: &types.ScriptTask{
				TaskProperties: types.TaskProperties{
					Type:           types.TaskTypeScript,
					CombinedOutput: true,
				},
				Interpreter: "python",
				Script:      `print("hello world")`,
			},
		},
	}
)

func fetchSchedule(token string, url string) (*types.Schedule, error) {
	return dummySchedule, nil
	// return nil, errors.New("unimplemented")
}

func fetchTasks(token string, url string) (map[string]*types.Task, error) {
	return dummyTasks, nil
	// return nil, errors.New("unimplemented")
}

func postResult(token string, url string, record *types.Record) error {
	if record == nil {
		return errors.New("nil record")
	}
	log.Println("posting record:", *record)
	return errors.New("unimplemented")
}
