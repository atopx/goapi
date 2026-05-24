package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"

	"goapi/common/logger"
	"goapi/common/system"

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
			resp.Code = system.ServerError
			resp.Message = "Internal Server Error"
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, resp)
		}()
		ctx.Next()
	}
}

func isBrokenPipe(rec any) bool {
	ne, ok := rec.(*net.OpError)
	if !ok {
		return false
	}
	se, ok := ne.Err.(*os.SyscallError)
	if !ok {
		return false
	}
	msg := strings.ToLower(se.Error())
	return strings.Contains(msg, "broken pipe") || strings.Contains(msg, "connection reset by peer")
}
