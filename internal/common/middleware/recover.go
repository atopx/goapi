package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"syscall"

	"goapi/internal/common/logger"
	"goapi/internal/common/system"

	"github.com/gin-gonic/gin"
)

// RecoverMiddleware 捕获 panic 并写入 ServerError 响应。
func RecoverMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			rec := recover()
			if rec == nil {
				return
			}

			brokenPipe := isBrokenPipe(rec)
			body, _ := httputil.DumpRequest(ctx.Request, false)

			if brokenPipe {
				logger.Error(ctx, "recovery from broken pipe",
					slog.String("body", string(body)),
					slog.Any("error", rec),
				)
				ctx.Abort()
				return
			}

			logger.Error(ctx, "recovery from panic",
				slog.Any("error", rec),
				slog.String("body", string(body)),
				slog.String("stack", string(debug.Stack())),
			)

			resp := system.GetResponse(ctx)
			if resp == nil {
				resp = system.NewResponse(ctx)
			}
			resp.Fail(system.ServerError, "Internal Server Error")
			// 与 Scheduler 保持一致：业务/系统错误统一返回 HTTP 200，错误语义由 resp.Code 承载。
			ctx.AbortWithStatusJSON(http.StatusOK, resp)
		}()
		ctx.Next()
	}
}

func isBrokenPipe(rec any) bool {
	err, ok := rec.(error)
	if !ok {
		return false
	}
	return errors.Is(err, syscall.EPIPE) || errors.Is(err, syscall.ECONNRESET)
}
