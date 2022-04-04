package schedule

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/robfig/cron"

	"github.com/LassiHeikkila/taskey/pkg/types"
)

type Executor interface {
	SetSchedule(schedule types.Schedule) error
	ConfigureTask(name string, task func()) error
	Start(ctx context.Context) error
	Stop() error
	Restart(ctx context.Context) error
}

func NewExecutor() (Executor, error) {
	e := &executor{
		ctx:          context.Background(),
		cronExecutor: cron.New(),
		scheduleConfig: scheduleConfig{
			cronConfig:       make(map[taskInstance]cron.Schedule),
			singleshotConfig: make(map[taskInstance]time.Time),
			periodicConfig:   make(map[taskInstance]time.Duration),
		},
		singleshotTimers: make(map[taskInstance]*time.Timer),
		periodicTickers:  make(map[taskInstance]*time.Ticker),
		tasks:            make(map[string]func()),
	}

	return e, nil
}

type executor struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	scheduleConfig scheduleConfig

	cronExecutor     *cron.Cron
	singleshotTimers map[taskInstance]*time.Timer
	periodicTickers  map[taskInstance]*time.Ticker

	scheduleChangeMutex sync.Mutex
	tasks               map[string]func()
}

type scheduleConfig struct {
	cronConfig       map[taskInstance]cron.Schedule
	singleshotConfig map[taskInstance]time.Time
	periodicConfig   map[taskInstance]time.Duration
}

// needed because a schedule may define same task
// to be run several times with periodic or singleshot timers
type taskInstance struct {
	name  string
	index int
}

func (e *executor) SetSchedule(schedule types.Schedule) error {
	e.scheduleChangeMutex.Lock()
	defer e.scheduleChangeMutex.Unlock()

	if len(schedule.CronTasks) > 0 {
		for _, ct := range schedule.CronTasks {
			if err := e.scheduleCronTask(ct); err != nil {
				return err
			}
		}
	}

	if len(schedule.PeriodicTasks) > 0 {
		for _, pt := range schedule.PeriodicTasks {
			if err := e.schedulePeriodicTask(pt); err != nil {
				return err
			}
		}
	}

	if len(schedule.SingleshotTasks) > 0 {
		for _, st := range schedule.SingleshotTasks {
			if err := e.scheduleSingleshotTask(st); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *executor) scheduleCronTask(task types.CronTask) error {
	cs, err := cron.Parse(task.When)
	if err != nil {
		return err
	}
	ti := taskInstance{
		name:  task.What,
		index: getTaskInstance(e.scheduleConfig.cronConfig, task.What),
	}
	e.scheduleConfig.cronConfig[ti] = cs
	return nil
}

func (e *executor) scheduleSingleshotTask(task types.SingleshotTask) error {
	ti := taskInstance{
		name:  task.What,
		index: getTaskInstance(e.scheduleConfig.singleshotConfig, task.What),
	}
	e.scheduleConfig.singleshotConfig[ti] = task.When
	return nil
}

func (e *executor) schedulePeriodicTask(task types.PeriodicTask) error {
	ti := taskInstance{
		name:  task.What,
		index: getTaskInstance(e.scheduleConfig.periodicConfig, task.What),
	}
	e.scheduleConfig.periodicConfig[ti] = task.Interval
	return nil
}

func (e *executor) ConfigureTask(name string, task func()) error {
	e.scheduleChangeMutex.Lock()
	defer e.scheduleChangeMutex.Unlock()

	e.tasks[name] = task
	return nil
}

func (e *executor) Start(ctx context.Context) error {
	e.ctx, e.ctxCancel = context.WithCancel(ctx)

	for k, v := range e.scheduleConfig.cronConfig {
		job := cron.FuncJob(func() {
			t, defined := e.tasks[k.name]
			if !defined {
				// TODO: log something?
				return
			}
			t()
		})
		e.cronExecutor.Schedule(v, job)
	}

	for k, v := range e.scheduleConfig.singleshotConfig {
		d := v.Sub(time.Now())
		t := time.NewTimer(d)
		e.singleshotTimers[k] = t
		go func() {
			select {
			case <-e.ctx.Done():
				return
			case <-t.C:
				t, defined := e.tasks[k.name]
				if !defined {
					// TODO: log something?
					return
				}
				t()
			}
		}()
	}

	for k, v := range e.scheduleConfig.periodicConfig {
		t := time.NewTicker(v)
		e.periodicTickers[k] = t
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-t.C:
					t, defined := e.tasks[k.name]
					if !defined {
						// TODO: log something?
						// TODO: continue instead of return to keep trying?
						return
					}
					t()
				}
			}
		}()
	}

	return nil
}

func (e *executor) Stop() error {
	if e.ctxCancel != nil {
		e.ctxCancel()
	}
	return errors.New("unimplemented")
}

func (e *executor) Restart(ctx context.Context) error {
	if err := e.Stop(); err != nil {
		return err
	}
	if err := e.Start(ctx); err != nil {
		return err
	}
	return nil
}

func getTaskInstance[V interface{}](m map[taskInstance]V, task string) int {
	i := 0
	for k := range m {
		if k.name == task {
			i++
		}
	}
	return i
}
