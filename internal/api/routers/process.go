package routers

import (
	handlers "LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"
	"github.com/gin-gonic/gin"
)

type ProcessRouter struct {
}

func (f *ProcessRouter) InitRouter(Router *gin.RouterGroup) {
	processRouter := Router.Group("process").Use(middleware.PasswordExpired())
	baseApi := handlers.ApiGroupApp.BaseApi
	{
		processRouter.GET("/ws", baseApi.ProcessWs)
		processRouter.POST("/stop", baseApi.StopProcess)
	}
}
