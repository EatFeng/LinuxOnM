package routers

import (
	handlers "LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"

	"github.com/gin-gonic/gin"
)

type LogRouter struct{}

func (s *LogRouter) InitRouter(Router *gin.RouterGroup) {
	operationRouter := Router.Group("log").Use(middleware.PasswordExpired())
	baseApi := handlers.ApiGroupApp.BaseApi
	{
		operationRouter.POST("/login", baseApi.GetLoginLog)
		operationRouter.POST("/operation", baseApi.GetOperationLog)
		operationRouter.GET("/system/files", baseApi.GetSystemFiles)
		operationRouter.POST("/ssh", baseApi.LoadSSHLog)
	}
}
