package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type HostRouter struct{}

func (s *HostRouter) InitRouter(Router *gin.RouterGroup) {
	hostRouter := Router.Group("host")
	baseApi := handler.ApiGroupApp.BaseApi
	{
		hostRouter.POST("", baseApi.CreateHost)
		hostRouter.POST("/search", baseApi.SearchHost)
		hostRouter.POST("/update", baseApi.UpdateHost)
		hostRouter.POST("/update/group", baseApi.UpdateHostGroup)
		hostRouter.POST("/tree", baseApi.HostTree)
		hostRouter.POST("/test/byid/:id", baseApi.TestByID)
		hostRouter.POST("/test/byinfo", baseApi.TestByInfo)
	}
}
