# goapi

> go1.26+

A modern Go REST API template — Gin · GORM · PostgreSQL · `log/slog` · gocron/v2.

## Features

- PostgreSQL-only via GORM (drop MySQL/SQLite, keep one well-tuned driver)
- TOML configuration via `pelletier/go-toml/v2`
- Standard-library structured logging (`log/slog`) with lumberjack rotation and `trace_id` propagation
- gocron/v2 scheduler with config-driven, reflection-based task registration (no hardcoded `switch` per task)
- Graceful shutdown on `SIGINT/SIGTERM/SIGQUIT`
- Lightweight `api/control/biz` layering
- Dockerfile (CGO-disabled multi-stage build)

## Installation

```bash
git clone https://github.com/atopx/goapi.git
cd goapi
cp conf/config.example.toml conf/config.toml   # edit DB credentials
go mod tidy
go run .
```

## Layout

```bash
goapi
├── common
│   ├── handle           # global db/redis handles
│   ├── logger           # slog + gorm logger adapter
│   ├── middleware       # gin middlewares (recover, request/response logging, trace)
│   ├── system           # api response envelope + status codes
│   └── utils            # sql helpers, pagination, range, trace id
├── conf                 # toml config loader and example files
├── docs                 # redoc / openapi (placeholder)
├── internal
│   ├── api              # thin route → control bridges
│   ├── biz              # business logic, organized per-action
│   ├── control          # generic request/response Controller
│   ├── model            # gorm models
│   ├── scheduler
│   │   ├── scheduler.go # gocron/v2 wrapper + reflection-based loadTasks
│   │   ├── tasks.go     # list of task constructors (only place to declare new tasks)
│   │   ├── task/        # base Task interface + factory + wrapper
│   │   └── health/      # sample task
│   └── server           # http.Server bootstrap + router
├── pkg                  # postgres / redis client factories
└── tests                # config tests
```

## Adding a new scheduled task

1. Create `internal/scheduler/<name>/task.go` with:
   ```go
   package <name>

   import "goapi/internal/scheduler/task"

   type Task struct { *task.BaseTask }
   func New(name string, args map[string]any) *Task {
       return &Task{BaseTask: task.New(name, args)}
   }
   func (t *Task) Run() error { /* ... */ return nil }
   ```
2. Append `<name>.New` to `internal/scheduler/tasks.go`.
3. Add a `[[scheduler]]` block in `conf/config.toml` with `name = "<name>"` and your cron spec.

The scheduler reflects on each constructor's return type, picks the last path segment of its package as the task name, and binds it to the matching config entry — no hardcoded `switch/case`.

## License

MIT
