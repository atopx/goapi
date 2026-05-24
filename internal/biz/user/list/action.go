package user_list

import (
	"fmt"
	"log/slog"
	"time"

	"goapi/common/logger"
	"goapi/common/utils"
)

func (c *Controller) Deal() (any, error) {
	params := c.Params.(*Params)

	fmt.Printf("%+v\n", params)

	reply := Reply{
		Total:    1,
		Filtered: 1,
		Page: utils.Pagination{
			Page: 1,
			Size: 1,
		},
		Records: []*Record{
			{
				Id:         1,
				Username:   "atopx",
				Nickname:   "atopx",
				Age:        20,
				CreateTime: time.Now().Format(time.DateTime),
				UpdateTime: time.Now().Format(time.DateTime),
			},
		},
	}

	logger.Info(c.Context(), "user_list success", slog.Int64("filtered", reply.Filtered))
	return reply, nil
}
