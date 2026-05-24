package server

import (
	"net/http"
	"time"

	"goapi/conf"

	"github.com/gin-gonic/gin"
)

// New 构造 *http.Server，由 main 负责生命周期管理（ListenAndServe + Shutdown）。
func New(cfg *conf.ServerConfig) *http.Server {
	app := gin.New()
	return &http.Server{
		Addr:           cfg.Addr,
		Handler:        router(app),
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
}
