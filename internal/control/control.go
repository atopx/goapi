package control

import (
	"log/slog"
	"net/http"

	"goapi/internal/common/logger"
	"goapi/internal/common/system"

	"github.com/gin-gonic/gin"
)

type Controller[P any] struct {
	context *gin.Context
	Params  *P
	err     error
}

func New[P any](ctx *gin.Context) *Controller[P] {
	params := new(P)
	err := ctx.ShouldBind(params)
	return &Controller[P]{
		context: ctx,
		Params:  params,
		err:     err,
	}
}

func (ctl *Controller[P]) Context() *gin.Context {
	return ctl.context
}

func (ctl *Controller[P]) Error() error {
	return ctl.err
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
		logger.Warn(ctx, "bind params error", slog.Any("error", err))
		resp.Fail(system.ClientError, err.Error())
	} else if data, err := ctl.Deal(); err != nil {
		resp.Fail(system.ServerError, err.Error())
	} else {
		resp.Success(data)
	}

	ctx.JSON(http.StatusOK, resp)
}
