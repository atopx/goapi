package task

import (
	"fmt"
	"reflect"
	"goapi/conf"
)

var factories = make(map[string]any)

func Register(name string, factory any) {
	factories[name] = factory
}

func Create(cfg conf.WorkerConfig) (Task, error) {
	if cfg.Disable {
		return nil, fmt.Errorf("disabled task")
	}
	t, ok := factories[cfg.Name]
	if !ok {
		return nil, fmt.Errorf("unimplemented task")
	}
	factory := reflect.ValueOf(t)
	args := []reflect.Value{
		reflect.ValueOf(cfg.Name),
		reflect.ValueOf(cfg.Args),
	}
	return factory.Call(args)[0].Interface().(Task), nil
}
