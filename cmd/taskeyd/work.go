package main

import (
	"context"
	"errors"

	"github.com/LassiHeikkila/taskey/pkg/schedule"
	"github.com/LassiHeikkila/taskey/pkg/types"
)

func executeSchedule(ctx context.Context, sched *types.Schedule, tasks map[string]*types.Task) error {
	if sched == nil {
		return errors.New("nil schedule")
	}
	if len(tasks) == 0 {
		return errors.New("no tasks defined")
	}
	executor, err := schedule.NewExecutor()
	if err != nil {
		return err
	}

	if err := executor.SetSchedule(*sched); err != nil {
		return err
	}

	execCb := taskExecCallback(func(name string, status int, output string) {

	})

	for name, task := range tasks {
		if err := executor.ConfigureTask(name, makeTask(task, execCb)); err != nil {
			return err
		}
	}

	err = executor.Start(ctx)
	defer executor.Stop()
	if err != nil {
		return err
	}
	<-ctx.Done()

	return nil
}