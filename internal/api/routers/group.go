package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type GroupRouter struct{}

func (a *GroupRouter) InitRouter(Router *gin.RouterGroup) {
	groupRouter := Router.Group("group")

	baseApi := handler.ApiGroupApp.BaseApi
	{
		groupRouter.POST("", baseApi.CreateGroup)
		groupRouter.POST("/del", baseApi.DeleteGroup)
		groupRouter.POST("/search", baseApi.ListGroup)
		groupRouter.POST("/update", baseApi.UpdateGroup)
	}
}
