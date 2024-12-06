package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type SettingRouter struct{}

func (s *SettingRouter) InitRouter(Router *gin.RouterGroup) {
	router := Router.Group("setting")

	baseApi := handler.ApiGroupApp.BaseApi
	{
		router.POST("/search", baseApi.GetSettingInfo)
	}
}
