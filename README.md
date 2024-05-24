# goapi

> go1.22+

This is a template for creating a REST API using Golang, Gin, Gorm, Zap, and Swagger.


## Features

- Support databases: MySQL, PostgreSQL, and SQLite databases
- Support docker
- Performance optimization
- Redoc and openapi documentation
- Full chain trace of logs (use uber-zap)
- Lightweight and fast
- Licensed under the MIT License

## Installation
1. Clone the repository:

  ```bash
  git clone https://github.com/atopx/goapi.git
  ```

2. Init your app
  
  ```bash
  cd goapi
  ./control init $NEW_NAME # example: ./control init newapp
  ```

3. Change the configuration file:

  ```bash
  cp conf/config.example.yaml conf/config.yaml
  ```

4. Run the application:

  ```bash
  go mod tidy
  go run main.go
  ```

5. Control script

  ```bash
  ./control { init $NEW_NAME | build | docker | docs | clean }
  ```

## Codes Layout

```bash
goapi
├── common # common package
│   ├── handle # global handlers
│   │   └── db.go
│   ├── logger # global logger
│   │   ├── gorm.go
│   │   └── logger.go
│   ├── middleware # router middleware
│   │   ├── context.go
│   │   └── recover.go
│   ├── system # api response
│   │   ├── consts.go
│   │   └── response.go
│   └── utils # frequently used functions
│       ├── sql.go
│       └── trace.go
├── conf # configuration file
│   ├── config.example.yaml
│   ├── config.go
│   └── config.yaml
├── docs # redoc and openapi documentation
├── internal # internal logic
│   ├── api # define api
│   │   ├── server.go
│   │   └── user.go
│   ├── biz # business logic
│   │   └── user
│   │       └── list
│   │           ├── action.go
│   │           └── schema.go
│   ├── control # api controller interface
│   │   └── control.go
│   ├── model # db model
│   │   └── user.go
│   ├── server # api server
│   │   ├── router.go
│   │   └── server.go
│   └── worker # async workers
│       ├── health
│       │   └── task.go
│       └── worker.go
├── pkg # third-party packages
│   ├── db.go
│   └── redis.go
└── tests # tests
    └── config_test.go
```
