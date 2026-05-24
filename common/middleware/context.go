package middleware

import (
	"log/slog"
	"time"

	"goapi/common/logger"
	"goapi/common/system"
	"goapi/common/utils"

	"github.com/gin-gonic/gin"
)

// ContextMiddleware 注入 traceId、Response 容器，并打印 request/response 日志。
func ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := system.NewResponse(ctx)

		traceId := ctx.Request.Header.Get(logger.TraceKey)
		if traceId == "" {
			traceId = utils.NewTraceId()
		}
		resp.TraceId = traceId
		ctx.Set(logger.TraceKey, resp.TraceId)

		logger.Info(ctx, "request",
			slog.String("method", ctx.Request.Method),
			slog.String("path", ctx.Request.URL.Path),
			slog.String("client", ctx.ClientIP()),
		)

		beginTime := time.Now()
		ctx.Next()

		elapsed := slog.String("elapsed", time.Since(beginTime).String())

		switch {
		case resp.Code < system.ClientError:
			logger.Info(ctx, "response", elapsed)
		case resp.Code >= system.ServerError:
			logger.Error(ctx, "response", elapsed, slog.String("system_error", resp.Message))
		default:
			logger.Warn(ctx, "response", elapsed, slog.String("client_error", resp.Message))
		}
	}
}
