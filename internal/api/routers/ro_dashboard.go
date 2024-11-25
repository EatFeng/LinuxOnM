package routers

import (
	"LinuxOnM/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

func (s *DashboardRouter) InitRouter(Router *gin.RouterGroup) {
	cmdRouter := Router.Group("dashboard")

	baseApi := handlers.ApiGroupApp.BaseApi
	{
		cmdRouter.GET("/base/os", baseApi.LoadDashboardOsInfo)
		cmdRouter.GET("/current/:ioOption/:netOption", baseApi.LoadDashboardCurrentInfo)
	}
}