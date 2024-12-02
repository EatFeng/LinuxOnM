package routers

import (
	handlers "LinuxOnM/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

type LogRouter struct{}

func (s *LogRouter) InitRouter(Router *gin.RouterGroup) {
	operationRouter := Router.Group("log")
	baseApi := handlers.ApiGroupApp.BaseApi
	{
		operationRouter.POST("/login", baseApi.GetLoginLog)
		operationRouter.POST("/operation", baseApi.GetOperationLog)
		operationRouter.GET("/system/files", baseApi.GetSystemFiles)
		operationRouter.POST("/ssh", baseApi.LoadSSHLog)
	}
}
