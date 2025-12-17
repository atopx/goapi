package task

import (
	"context"
	"fmt"
	"goapi/common/logger"
	"goapi/common/trace"
	"time"

	"go.uber.org/zap"
)

type Task interface {
	Run() error
	GetName() string
	Context() context.Context
}

type BaseTask struct {
	Ctx  context.Context
	Name string
	Args map[string]any
}

func New(name string, args map[string]any) *BaseTask {
	return &BaseTask{
		Ctx:  trace.NewDefaultTraceContext(),
		Name: name,
		Args: args,
	}
}

func (t *BaseTask) GetName() string {
	return t.Name
}

func (t *BaseTask) Context() context.Context {
	return t.Ctx
}

func (t *BaseTask) Start() error {
	return fmt.Errorf("unimplemented task")
}

func Wrapper(task Task) func() error {
	return func() error {
		ctx := task.Context()
		name := zap.String("task", task.GetName())
		logger.Info(ctx, "task start", name)
		start := time.Now()
		if err := task.Run(); err != nil {
			logger.Error(task.Context(), "task failed", name, zap.Error(err))
		}
		logger.Info(ctx, "task finish", name, zap.String("elapsed", time.Since(start).String()))
		return nil
	}
}
