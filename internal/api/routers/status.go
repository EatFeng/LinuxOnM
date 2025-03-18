package routers

import (
	"LinuxOnM/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

type StatusRouter struct{}

func (s *StatusRouter) InitRouter(Router *gin.RouterGroup) {
	runRouter := Router.Group("status")

	baseApi := handlers.ApiGroupApp.BaseApi
	{
		runRouter.GET("/current", baseApi.GetCurrentStatus)
	}
}
