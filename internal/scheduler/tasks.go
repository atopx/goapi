package scheduler

import "goapi/internal/scheduler/health"

// tasks 列出所有可用的任务构造函数。
// 新增任务时：
//  1. 在 internal/scheduler/<name>/ 下创建包，提供 func New(name string, args map[string]any) *Task。
//  2. 在这里 append 构造函数引用。
//  3. 在 config.toml 添加 [[scheduler]] 项，name 与包名一致即可被反射注册。
//
// 运行期不再有 switch/case 硬编码任务名。
var tasks = []any{
	health.New,
}
