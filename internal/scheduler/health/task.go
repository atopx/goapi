package health

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"goapi/conf"
	"goapi/internal/common/logger"
	"goapi/internal/scheduler/task"
)

// Task 是示例健康检查任务，生产环境请按需调整或删除。
type Task struct {
	*task.BaseTask
}

func New(name string, args map[string]any) *Task {
	return &Task{BaseTask: task.New(name, args)}
}

func (t *Task) Run() error {
	ctx := t.Context()
	url := fmt.Sprintf("http://%s/server/health", conf.Get().Server.Addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	if v, ok := ctx.Value(logger.TraceKey).(string); ok && v != "" {
		req.Header.Set(logger.TraceKey, v)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(ctx, "health request failed", slog.Any("error", err))
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	logger.Info(ctx, "health ok",
		slog.Int("status", resp.StatusCode),
		slog.String("body", string(body)),
	)
	return nil
}
