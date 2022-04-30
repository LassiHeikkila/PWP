package main

import (
	"bytes"
	stdjson "encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func setAuthorizationHeader(req *http.Request, token string) {
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", token))
}

func doGetRequest(req *http.Request, v any) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	dec := stdjson.NewDecoder(resp.Body)

	err = dec.Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func fetchSchedule(token string, url string, org string) (*types.Schedule, error) {
	if token == "" && url == "" {
		return dummySchedule, nil
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/api/v1/%s/machines/self/schedule/",
			url, org,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, token)

	type scheduleResponse struct {
		Code    int             `json:"code"`
		Message string          `json:"msg"`
		Payload *types.Schedule `json:"payload"`
	}

	var resp scheduleResponse

	err = doGetRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New("no schedule found")
	}

	return resp.Payload, nil
}

func fetchTasks(token string, url string, org string) (map[string]*types.Task, error) {
	if token == "" && url == "" {
		return dummyTasks, nil
	}

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf(
			"%s/api/v1/%s/machines/self/tasks/",
			url, org,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}
	setAuthorizationHeader(req, token)

	type tasksResponse struct {
		Code    int                      `json:"code"`
		Message string                   `json:"msg"`
		Payload []map[string]interface{} `json:"payload"`
	}

	var resp tasksResponse

	err = doGetRequest(req, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != http.StatusOK {
		return nil, errors.New("no tasks found")
	}

	m := make(map[string]*types.Task, len(resp.Payload))
	for _, t := range resp.Payload {
		task := unmarshalTask(t)
		m[task.Name] = task
	}

	return m, nil
}

func getValue[V any](m map[string]interface{}, key string) V {
	val, ok := m[key]
	if !ok {
		var zero V
		return zero
	}
	v, ok := val.(V)
	if !ok {
		var zero V
		return zero
	}
	return v
}

func getSlice[V any](m map[string]interface{}, key string) []V {
	val, ok := m[key]
	if !ok {
		var zero []V
		return zero
	}
	v, ok := val.([]interface{})
	if !ok {
		var zero []V
		return zero
	}

	r := make([]V, 0, len(v))
	for _, item := range v {
		r = append(r, item.(V))
	}

	return r
}

func unmarshalTask(v map[string]interface{}) *types.Task {
	log.Println(v)
	// TODO: fix unsafe unmarshalling
	name := getValue[string](v, "name")
	description := getValue[string](v, "description")

	c := getValue[map[string]interface{}](v, "content")

	tp := getValue[string](c, "type")
	combo := getValue[bool](c, "combinedOutput")

	var content any
	switch tp {
	case types.TaskTypeCmd:
		program := getValue[string](c, "program")
		args := getSlice[string](c, "args")
		content = &types.CmdTask{
			TaskProperties: types.TaskProperties{
				Type:           tp,
				CombinedOutput: combo,
			},
			Program: program,
			Args:    args,
		}
	case types.TaskTypeScript:
		interpreter := getValue[string](c, "interpreter")
		scriptBody := getValue[string](c, "script")
		content = &types.ScriptTask{
			TaskProperties: types.TaskProperties{
				Type:           tp,
				CombinedOutput: combo,
			},
			Interpreter: interpreter,
			Script:      scriptBody,
		}
	}

	return &types.Task{
		Name:        name,
		Description: description,
		Content:     content,
	}
}

func postResult(token string, url string, org string, record *types.Record) error {
	if record == nil {
		return errors.New("nil record")
	}

	if token == "" && url == "" {
		log.Println("posting record:", *record)
		return nil
	}

	body := bytes.Buffer{}
	enc := stdjson.NewEncoder(&body)

	err := enc.Encode(record)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"%s/api/v1/%s/machines/self/records/",
			url, org,
		),
		&body,
	)
	if err != nil {
		return err
	}
	setAuthorizationHeader(req, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type response struct {
		Code    int    `json:"code"`
		Message string `json:"msg"`
	}

	dec := stdjson.NewDecoder(resp.Body)
	var v response
	err = dec.Decode(&v)
	if err != nil {
		return err
	}

	if v.Code != http.StatusOK {
		return fmt.Errorf("non-ok response: %d", v.Code)
	}

	return nil
}
