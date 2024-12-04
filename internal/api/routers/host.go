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
		hostRouter.POST("/del", baseApi.DeleteHost)
		hostRouter.POST("/search", baseApi.SearchHost)
		hostRouter.POST("/update", baseApi.UpdateHost)
		hostRouter.POST("/update/group", baseApi.UpdateHostGroup)
		hostRouter.POST("/tree", baseApi.HostTree)
		hostRouter.POST("/test/byid/:id", baseApi.TestByID)
		hostRouter.POST("/test/byinfo", baseApi.TestByInfo)

		hostRouter.GET("/command", baseApi.ListCommand)
		hostRouter.POST("/command", baseApi.CreateCommand)
		hostRouter.POST("/command/del", baseApi.DeleteCommand)
		hostRouter.POST("/command/update", baseApi.UpdateCommand)
		hostRouter.POST("/command/search", baseApi.SearchCommand)
		hostRouter.GET("/command/tree", baseApi.SearchCommandTree)
	}
}
