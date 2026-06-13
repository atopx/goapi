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
		var sb strings.Builder
		sb.WriteString(field)
		sb.WriteString(" between ? and ?")
		return tx.Where(sb.String(), r.Start, r.End)
	}
}

type Pagination struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// RecordsCap 计算当前页面预期返回的记录数（用于预分配切片容量）。
// total 是过滤后的总数，已知页码和页大小后，当前页最多承载 min(size, total - (page-1)*size) 条。
func (p *Pagination) RecordsCap(total int64) int {
	page := p.Page
	if page <= 0 {
		page = 1
	}
	size := p.Size
	if size <= 0 {
		size = 10
	}
	offset := int64(page-1) * int64(size)
	remaining := total - offset
	if remaining <= 0 {
		return 0
	}
	if remaining > int64(size) {
		return size
	}
	return int(remaining)
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
	var sb strings.Builder
	sb.WriteByte('%')
	sb.WriteString(key)
	sb.WriteByte('%')
	return sb.String()
}
