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
		settingRouter.POST("/update", baseApi.UpdateSetting)
		settingRouter.GET("/interface", baseApi.LoadInterfaceAddr)
	}
}
