package health

import (
	"context"
	"fmt"
	"goapi/common/logger"
	"goapi/common/utils"
	"goapi/conf"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// Health 这里只是一个示例任务，生产环境请禁用或删除
func Health() {
	traceId := utils.NewTraceId()
	ctx := context.WithValue(context.Background(), logger.TraceKey(), traceId)
	url := fmt.Sprintf("http://%s/api/server/health", conf.Get().Server.Addr)
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("traceId", traceId)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		logger.Error(ctx, "health", zap.Error(err))
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	logger.Info(ctx, "health", zap.Int("status", resp.StatusCode), zap.ByteString("resp", body))
}
