package userlist

import (
	"fmt"
	"log/slog"
	"time"

	"goapi/internal/common/handle"
	"goapi/internal/common/logger"
	"goapi/internal/common/utils"
	"goapi/internal/model"
)

func (c *Controller) Deal() (any, error) {
	tx := handle.DB(c.Context()).Model(&model.User{})

	var reply Reply

	if err := tx.Count(&reply.Total).Error; err != nil {
		return nil, fmt.Errorf("user_list total count: %w", err)
	}

	if c.Params.Keyword != "" {
		key := utils.Like(c.Params.Keyword)
		tx = tx.Where("username like ? or nickname like ?", key, key)
	}

	if c.Params.AgeRange.IsValid() {
		tx = tx.Scopes(c.Params.AgeRange.Between("age"))
	}

	if c.Params.CreateTimeRange.IsValid() {
		tx = tx.Scopes(c.Params.CreateTimeRange.Between("created_at"))
	}

	if err := tx.Count(&reply.Filtered).Error; err != nil {
		return nil, fmt.Errorf("user_list filtered count: %w", err)
	}

	users := make([]*model.User, 0, c.Params.Page.RecordsCap(reply.Filtered))

	if err := tx.Scopes(c.Params.Page.Paging).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("user_list find: %w", err)
	}

	reply.Page = c.Params.Page
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
