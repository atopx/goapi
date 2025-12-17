package middleware

import (
	"goapi/common/system"
	"goapi/common/trace"
	"time"

	"goapi/common/logger"

	"github.com/atopx/clever/general"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ContextMiddleware 日志中间件
func ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp := system.NewResponse(ctx)

		// 如果请求头里携带了traceId, 则不再重新生成
		traceId := ctx.Request.Header.Get("trace_id")
		if traceId == general.Empty {
			traceId = trace.NewTraceID()
		}
		resp.TraceId = traceId
		// 将 trace 写入上下文供日志获取
		ctx.Set(string(trace.ContextTrace), resp.TraceId)

		// request log
		logger.Info(ctx, "request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("client", ctx.ClientIP()),
		)

		// 执行接口逻辑
		beginTime := time.Now()
		ctx.Next()

		elapsed := zap.String("elapsed", time.Since(beginTime).String())

		if resp.Code < system.ClientError {
			logger.Info(ctx, "response", elapsed)
		} else if resp.Code >= system.ServerError {
			logger.Error(ctx, "response", elapsed, zap.String("systemError", resp.Message))
		} else {
			logger.Warn(ctx, "response", elapsed, zap.String("clientError", resp.Message))
		}
	}
}
