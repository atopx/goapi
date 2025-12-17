package health

import (
	"fmt"
	"goapi/common/logger"
	"goapi/common/trace"
	"goapi/internal/app"
	"goapi/internal/scheduler/task"
	"io"
	"net/http"

	"go.uber.org/zap"
)

type Task struct {
	*task.BaseTask
}

func New(name string, args map[string]any) *Task {
	return &Task{BaseTask: task.New(name, args)}
}

func (t *Task) Run() error {
	traceId := trace.NewTraceID()
	ctx := trace.WithTraceID(t.Ctx, traceId)
	url := fmt.Sprintf("http://%s/server/health", app.Infra().Config().Server.Addr)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("trace_id", traceId)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Error(ctx, "health", zap.Error(err))
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	logger.Info(ctx, "health", zap.Int("status", resp.StatusCode), zap.ByteString("resp", body))
	return nil
}
