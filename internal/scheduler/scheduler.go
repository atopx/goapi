package scheduler

import (
	"context"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"goapi/common/logger"
	"goapi/conf"
	"goapi/internal/scheduler/task"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

type Scheduler struct {
	ctx       context.Context
	scheduler gocron.Scheduler
	jobs      map[string]gocron.Job
}

func loadTasks() {
	for _, t := range tasks {
		name := strings.TrimPrefix(reflect.TypeOf(t).Out(0).String(), "*")
		task.Register(name, t)
	}
}

func New(ctx context.Context, cfgs []conf.WorkerConfig) *Scheduler {
	loadTasks()
	s, _ := gocron.NewScheduler()
	scheduler := &Scheduler{
		ctx:       ctx,
		scheduler: s,
		jobs:      make(map[string]gocron.Job),
	}
	for _, cfg := range cfgs {
		t, err := task.Create(cfg)
		if err != nil {
			logger.Warn(ctx, "init task failed", zap.String("task", cfg.Name), zap.Error(err))
			continue
		}
		job, err := scheduler.scheduler.NewJob(
			gocron.CronJob(cfg.Spec, true),
			gocron.NewTask(task.Wrapper(t)),
			gocron.WithName(cfg.Name),
			gocron.WithSingletonMode(gocron.LimitModeWait),
		)
		if err != nil {
			logger.Error(ctx, "register task error", zap.Error(err), zap.String("task", cfg.Name))
			continue
		}
		scheduler.jobs[cfg.Name] = job
	}
	return scheduler
}

func (s *Scheduler) Start() {
	s.scheduler.Start()
}

func (s *Scheduler) Runloop() {
	s.scheduler.Start()
	ctx, stop := signal.NotifyContext(s.ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	_ = s.Shutdown()
}

func (s *Scheduler) Shutdown() error {
	return s.scheduler.Shutdown()
}
