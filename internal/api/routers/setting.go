package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type SettingRouter struct{}

func (s *SettingRouter) InitRouter(Router *gin.RouterGroup) {
	router := Router.Group("setting")
	settingRouter := Router.Group("setting")
	baseApi := handler.ApiGroupApp.BaseApi
	{
		router.POST("/search", baseApi.GetSettingInfo)
		router.POST("/expired/handle", baseApi.HandlePasswordExpired)
		settingRouter.GET("/search/available", baseApi.GetSystemAvailable)
		settingRouter.GET("/interface", baseApi.LoadInterfaceAddr)
		settingRouter.POST("/update", baseApi.UpdateSetting)
		settingRouter.POST("/update/password", baseApi.UpdatePassword)
		settingRouter.POST("/update/proxy", baseApi.UpdateProxy)
		settingRouter.POST("/update/bind", baseApi.UpdateBindInfo)
		settingRouter.POST("/update/port", baseApi.UpdatePort)
		settingRouter.POST("/ssl/update", baseApi.UpdateSSL)
	}
}
