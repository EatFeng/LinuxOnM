package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type CronjobRouter struct{}

func (s *CronjobRouter) InitRouter(Router *gin.RouterGroup) {
	cmdRouter := Router.Group("cronjob")
	baseApi := handler.ApiGroupApp.BaseApi
	{
		cmdRouter.POST("", baseApi.CreateCronjob)
	}
}
