package user_list

import (
	"fmt"
	"goapi/common/logger"
	"goapi/common/utils"
	"goapi/internal/app"
	"goapi/internal/model"
	"time"

	"github.com/atopx/clever/general"
)

func (c *Controller) Deal() (any, error) {
	params := c.Params.(*Params)

	tx := app.Infra().DB(c.Context()).Model(&model.User{})

	var reply Reply

	if err := tx.Count(&reply.Total).Error; err != nil {
		return nil, fmt.Errorf("user_list total count error: %s", err)
	}

	if params.Keyword != general.Empty {
		key := utils.Like(params.Keyword)
		tx.Where("username like ? or nickname like ?", key, key)
	}

	if params.AgeRange.IsValid() {
		tx.Scopes(params.AgeRange.Between("age"))
	}

	if params.CreateTimeRange.IsValid() {
		tx.Scopes(params.AgeRange.Between("create_at"))
	}

	if err := tx.Count(&reply.Filtered).Error; err != nil {
		return nil, fmt.Errorf("user_list filtered count error: %s", err)
	}

	users := make([]*model.User, 0, params.Page.RecordsCap(reply.Total))

	if err := tx.Scopes(params.Page.Paging).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("user_list find error: %s", err)
	}

	reply.Page = params.Page
	reply.Records = make([]*Record, 0, len(users))

	for _, user := range users {
		record := &Record{
			Id:         user.ID,
			Username:   user.Username,
			Nickname:   user.Nickname,
			Age:        user.Age,
			UpdateTime: user.CreatedAt.Format(time.DateTime),
			CreateTime: user.UpdatedAt.Format(time.DateTime),
		}
		reply.Records = append(reply.Records, record)
	}
	logger.Info(c.Context(), "success")
	return reply, nil
}
