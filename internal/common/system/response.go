package system

import (
	"github.com/gin-gonic/gin"
)

// Response 通用Response结构
// Code:
//   - [20000000, 40000000) success
//   - [40000000, 50000000) client error
//   - [50000000, 99999999] system error
type Response struct {
	TraceId string `json:"traceId"`
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Success 写入成功响应。
func (r *Response) Success(data any) {
	r.Code = SuccessCode
	r.Message = "OK"
	r.Data = data
}

// Fail 写入失败响应，code 取 ClientError / ServerError。
func (r *Response) Fail(code Code, message string) {
	r.Code = code
	r.Message = message
}

func (r *Response) IsSuccess() bool     { return r.Code < ClientError }
func (r *Response) IsClientError() bool { return r.Code >= ClientError && r.Code < ServerError }
func (r *Response) IsServerError() bool { return r.Code >= ServerError }

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
