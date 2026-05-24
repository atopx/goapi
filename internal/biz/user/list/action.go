package user_list

import (
	"fmt"
	"log/slog"
	"time"

	"goapi/common/handle"
	"goapi/common/logger"
	"goapi/common/utils"
	"goapi/internal/model"
)

func (c *Controller) Deal() (any, error) {
	params := c.Params.(*Params)

	tx := handle.DB(c.Context()).Model(&model.User{})

	var reply Reply

	if err := tx.Count(&reply.Total).Error; err != nil {
		return nil, fmt.Errorf("user_list total count: %w", err)
	}

	if params.Keyword != "" {
		key := utils.Like(params.Keyword)
		tx = tx.Where("username like ? or nickname like ?", key, key)
	}

	if params.AgeRange.IsValid() {
		tx = tx.Scopes(params.AgeRange.Between("age"))
	}

	if params.CreateTimeRange.IsValid() {
		tx = tx.Scopes(params.CreateTimeRange.Between("created_at"))
	}

	if err := tx.Count(&reply.Filtered).Error; err != nil {
		return nil, fmt.Errorf("user_list filtered count: %w", err)
	}

	users := make([]*model.User, 0, params.Page.RecordsCap(reply.Filtered))

	if err := tx.Scopes(params.Page.Paging).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("user_list find: %w", err)
	}

	reply.Page = params.Page
	reply.Records = make([]*Record, 0, len(users))

	for _, user := range users {
		reply.Records = append(reply.Records, &Record{
			Id:         user.ID,
			Username:   user.Username,
			Nickname:   user.Nickname,
			Age:        user.Age,
			CreateTime: user.CreatedAt.Format(time.DateTime),
			UpdateTime: user.UpdatedAt.Format(time.DateTime),
		})
	}
	logger.Info(c.Context(), "user_list success", slog.Int64("filtered", reply.Filtered))
	return reply, nil
}
