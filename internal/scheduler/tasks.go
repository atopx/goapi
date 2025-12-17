package scheduler

import (
	"goapi/internal/scheduler/health"
)

var tasks = []any{
	health.New,
}
