package control

import (
	"errors"
	"goapi/common/logger"
	"goapi/common/system"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	context *gin.Context
	Params  any
	err     error
}

func New(ctx *gin.Context, params any) *Controller {
	err := ctx.ShouldBind(params)

	return &Controller{
		context: ctx,
		Params:  params,
		err:     err,
	}
}

func (ctl *Controller) Context() *gin.Context {
	return ctl.context
}

func (ctl *Controller) Error() error {
	return ctl.err
}

func (ctl *Controller) Deal() (any, error) {
	return nil, errors.New("unimplemented")
}

type Handler interface {
	Context() *gin.Context
	Deal() (any, error)
	Error() error
}

func Scheduler(ctl Handler) {
	ctx := ctl.Context()
	resp := system.GetResponse(ctx)

	if err := ctl.Error(); err != nil {
		logger.Warn(ctx, "bind params error", zap.Error(err))
		resp.Code = system.ClientError
		resp.Message = err.Error()
	} else if resp.Data, err = ctl.Deal(); err != nil {
		resp.Code = system.ServerError
		resp.Message = err.Error()
	} else {
		resp.Message = "OK"
	}

	ctx.JSON(http.StatusOK, resp)
}
