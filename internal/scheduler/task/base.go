package task

import (
	"context"
	"log/slog"
	"time"

	"goapi/common/logger"
	"goapi/common/utils"
)

// Task 是所有调度任务必须实现的接口。
type Task interface {
	Run() error
	GetName() string
	Context() context.Context
}

// BaseTask 为任务提供默认的 ctx/name/args 字段，业务任务通过嵌入复用。
type BaseTask struct {
	Ctx  context.Context
	Name string
	Args map[string]any
}

// New 创建带 trace_id 的 BaseTask。每次调用生成新的 traceId，便于按任务维度串联日志。
func New(name string, args map[string]any) *BaseTask {
	ctx := context.WithValue(context.Background(), logger.TraceKey, utils.NewTraceId())
	return &BaseTask{Ctx: ctx, Name: name, Args: args}
}

func (t *BaseTask) GetName() string          { return t.Name }
func (t *BaseTask) Context() context.Context { return t.Ctx }

// Wrapper 将 Task 包装成 gocron 可调度的 func()，统一记录开始/结束/耗时/错误。
func Wrapper(t Task) func() {
	return func() {
		ctx := t.Context()
		nameAttr := slog.String("task", t.GetName())
		start := time.Now()
		logger.Info(ctx, "task start", nameAttr)
		if err := t.Run(); err != nil {
			logger.Error(ctx, "task failed", nameAttr, slog.Any("error", err))
		}
		logger.Info(ctx, "task finish", nameAttr, slog.String("elapsed", time.Since(start).String()))
	}
}
