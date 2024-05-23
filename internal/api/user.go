package api

import (
	user_list "goapi/internal/biz/user/list"
	"goapi/internal/control"

	"github.com/gin-gonic/gin"
)

// UserList
// @Summary 用户列表
// @Description 用户列表
// @Tags 用户
// @Accept json
// @Produce json
// @Param param body user_list.Params true "用户列表参数"
// @Response 200 object system.Response{data=user_list.Reply} "调用成功"
// @Router /api/v1/user/list [post]
func UserList(ctx *gin.Context) {

	control.Scheduler(user_list.NewController(ctx))
}
