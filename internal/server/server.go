package server

import (
	"goapi/conf"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Start(cfg *conf.ServerConfig, errors chan<- error) {
	app := gin.New()
	srv := &http.Server{
		Addr:           cfg.Addr,
		Handler:        router(app),
		ReadTimeout:    time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}
	go func() {
		log.Printf("start server http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			errors <- err
		}
	}()
}
