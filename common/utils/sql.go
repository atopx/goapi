package utils

import (
	"strings"

	"gorm.io/gorm"
)

type Range struct {
	Start int64 `json:"left"`
	End   int64 `json:"right"`
}

func (r *Range) IsValid() bool {
	return r.End > r.Start
}

func (r *Range) Between(field string) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		var sql strings.Builder
		sql.WriteString(field)
		sql.WriteString(" between ? and ?")
		return tx.Where(sql, r.Start, r.End)
	}
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

func (p *Pagination) RecordsCap(num int64) int {
	cap := p.Page*p.Size - int(num)
	if cap < 0 {
		return 0
	}
	if cap > p.Size {
		cap = p.Size
	}
	return cap
}

func (p *Pagination) Paging(db *gorm.DB) *gorm.DB {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Size <= 0 {
		p.Size = 10
	}
	return db.Offset((p.Page - 1) * p.Size).Limit(p.Size)
}

func Like(key string) string {
	var builder strings.Builder
	builder.WriteByte('%')
	builder.WriteString(key)
	builder.WriteByte('%')
	return builder.String()
}
