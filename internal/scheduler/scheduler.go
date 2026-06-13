package scheduler

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"reflect"

	"goapi/conf"
	"goapi/internal/common/logger"
	"goapi/internal/scheduler/task"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	ctx       context.Context
	scheduler gocron.Scheduler
	jobs      map[string]gocron.Job
}

// loadTasks 反射 tasks 切片中每个构造函数返回类型所在 package 的末段作为任务名注册。
// 例如 health.New 返回 *health.Task，注册名为 "health"。
// 这样配置 [[scheduler]] name = "health" 即可命中，不需要硬编码任务名。
func loadTasks() {
	for _, t := range tasks {
		out := reflect.TypeOf(t).Out(0)
		if out.Kind() == reflect.Ptr {
			out = out.Elem()
		}
		name := path.Base(out.PkgPath())
		task.Register(name, t)
	}
}

func New(ctx context.Context, cfgs []conf.WorkerConfig) (*Scheduler, error) {
	loadTasks()
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("create gocron scheduler: %w", err)
	}
	sch := &Scheduler{ctx: ctx, scheduler: s, jobs: make(map[string]gocron.Job)}
	for _, cfg := range cfgs {
		t, err := task.Create(cfg)
		if err != nil {
			logger.Warn(ctx, "init task skipped",
				slog.String("task", cfg.Name), slog.Any("error", err))
			continue
		}
		job, err := sch.scheduler.NewJob(
			gocron.CronJob(cfg.Spec, true),
			gocron.NewTask(task.Wrapper(t)),
			gocron.WithName(cfg.Name),
			gocron.WithSingletonMode(gocron.LimitModeWait),
		)
		if err != nil {
			logger.Error(ctx, "register task failed",
				slog.String("task", cfg.Name), slog.Any("error", err))
			continue
		}
		sch.jobs[cfg.Name] = job
		logger.Info(ctx, "task registered",
			slog.String("task", cfg.Name), slog.String("spec", cfg.Spec))
	}
	return sch, nil
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}
