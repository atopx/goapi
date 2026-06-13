package userlist

import (
	"goapi/internal/common/utils"
	"goapi/internal/control"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	*control.Controller[Params]
}

func NewController(ctx *gin.Context) *Controller {
	return &Controller{control.New[Params](ctx)}
}

type Params struct {
	Page            utils.Pagination `json:"page"`
	Keyword         string           `json:"keyword"`
	AgeRange        utils.Range      `json:"ageRange"`
	CreateTimeRange utils.Range      `json:"createTimeRange"`
}

type Reply struct {
	Total    int64            `json:"total"`
	Filtered int64            `json:"filtered"`
	Page     utils.Pagination `json:"page"`
	Records  []*Record        `json:"data"`
}

type Record struct {
	Id         uint   `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Age        int    `json:"age"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}
