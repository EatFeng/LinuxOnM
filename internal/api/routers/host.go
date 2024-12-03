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
		hostRouter.POST("/test/byid/:id", baseApi.TestByID)
	}
}
