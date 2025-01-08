package routers

import (
	"LinuxOnM/internal/api/handlers"
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
		processRouter.POST("/kill", baseApi.KillProcess)
		processRouter.POST("/content", baseApi.GetProcessContent)
		processRouter.POST("/start", baseApi.StartProcess)
		processRouter.POST("/stop", baseApi.StopProcess)
		processRouter.POST("/enable", baseApi.EnableProcess)
		processRouter.POST("/disable", baseApi.DisableProcess)
		processRouter.POST("/status", baseApi.StatusProcess)
		processRouter.POST("", baseApi.CreateProcess)
	}
}
