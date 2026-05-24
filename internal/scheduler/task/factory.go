package task

import (
	"fmt"
	"reflect"

	"goapi/conf"
)

// 任务工厂注册表。Key 为 conf.WorkerConfig.Name，Value 为 func(name, args) *XxxTask。
var factories = make(map[string]any)

// Register 由 scheduler 包在初始化阶段反射调用，把任务构造函数挂入注册表。
func Register(name string, factory any) {
	factories[name] = factory
}

// Create 按 cfg.Name 找到构造函数并通过反射调用，返回构造好的 Task 实例。
// 反射签名约定：func(name string, args map[string]any) T，T 必须实现 Task。
func Create(cfg conf.WorkerConfig) (Task, error) {
	if cfg.Disable {
		return nil, fmt.Errorf("task disabled: %s", cfg.Name)
	}
	f, ok := factories[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("unregistered task: %s", cfg.Name)
	}
	fv := reflect.ValueOf(f)
	if fv.Kind() != reflect.Func {
		return nil, fmt.Errorf("factory for %s is not a function", cfg.Name)
	}
	out := fv.Call([]reflect.Value{
		reflect.ValueOf(cfg.Name),
		reflect.ValueOf(cfg.Args),
	})
	t, ok := out[0].Interface().(Task)
	if !ok {
		return nil, fmt.Errorf("factory for %s did not return Task", cfg.Name)
	}
	return t, nil
}
