package main

import (
	"errors"

	"github.com/LassiHeikkila/taskey/pkg/types"
)

func fetchSchedule(token string, url string) (*types.Schedule, error) {
	return nil, errors.New("unimplemented")
}

func fetchTasks(token string, url string) (map[string]*types.Task, error) {
	return nil, errors.New("unimplemented")
}

func postResult(token string, url string, status int, output string) error {
	return errors.New("unimplemented")
}
