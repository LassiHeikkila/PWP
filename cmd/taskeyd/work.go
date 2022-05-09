package main

import (
	"context"
	"errors"
	"log"
	"time"

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
		log.Println("executed task", name, "with status", status) //, "and output:\n", output)
		rec := types.Record{
			TaskName:   name,
			ExecutedAt: time.Now(),
			Status:     status,
			Output:     output,
		}
		if err := postResult(config.AccessToken, config.URL, config.Organization, &rec); err != nil {
			log.Println("error posting result:", err)
		}
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
