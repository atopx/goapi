package system

import (
	"github.com/gin-gonic/gin"
)

// Response 通用Response结构
// Code:
//   - [20000000, 40000000] success
//   - [40000000, 50000000] client error
//   - [50000000, 99999999] system error
type Response struct {
	TraceId string `json:"traceId"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func NewResponse(ctx *gin.Context) *Response {
	resp := &Response{Code: SuccessCode}
	ctx.Set(ResponseKey, resp)
	return resp
}

func GetResponse(ctx *gin.Context) *Response {
	if value, ok := ctx.Get(ResponseKey); ok {
		return value.(*Response)
	}
	return nil
}
