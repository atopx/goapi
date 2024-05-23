package server

import (
	"goapi/common/middleware"
	"goapi/internal/api"

	"github.com/gin-gonic/gin"
)

func router(app *gin.Engine) *gin.Engine {

	// middleware
	app.Use(middleware.RecoverMiddleware())
	app.Use(middleware.ContextMiddleware())

	// server router
	serverGroup := app.Group("/server")
	{
		serverGroup.GET("/health", api.ServerHealth)
		serverGroup.GET("/config", api.ServerConfig)
	}

	// v1 router
	v1Group := app.Group("/api/v1")
	{
		// user apis
		userGroupV1 := v1Group.Group("/user")
		{
			userGroupV1.POST("/list", api.UserList)
		}

	}

	return app
}
