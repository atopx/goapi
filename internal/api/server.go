package api

import (
	"goapi/common/system"
	"goapi/internal/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerHealth
// @summary 服务检测
// @Tags 系统
// @Response 200 object system.Response "调用成功"
// @Router /server/health [get]
func ServerHealth(ctx *gin.Context) {
	resp := system.GetResponse(ctx)
	ctx.JSON(http.StatusOK, resp)
}

// ServerConfig
// @summary 服务配置
// @Tags 系统
// @Response 200 object system.Response{data=conf.Config} "调用成功"
// @Router /server/config [get]
func ServerConfig(ctx *gin.Context) {
	resp := system.GetResponse(ctx)
	resp.Data = app.Infra().Config
	ctx.JSON(http.StatusOK, resp)
}
