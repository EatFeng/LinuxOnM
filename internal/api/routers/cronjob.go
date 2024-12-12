package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"
	"github.com/gin-gonic/gin"
)

type CronjobRouter struct{}

func (s *CronjobRouter) InitRouter(Router *gin.RouterGroup) {
	cmdRouter := Router.Group("cronjob").Use(middleware.PasswordExpired())
	baseApi := handler.ApiGroupApp.BaseApi
	{
		cmdRouter.POST("", baseApi.CreateCronjob)
		cmdRouter.POST("/del", baseApi.DeleteCronjob)
		cmdRouter.POST("/search", baseApi.SearchCronjob)
		cmdRouter.POST("/update", baseApi.UpdateCronjob)
		cmdRouter.POST("/status", baseApi.UpdateCronjobStatus)
		cmdRouter.POST("/handle", baseApi.HandleOnce)
		cmdRouter.POST("/record/search", baseApi.SearchJobRecords)
		cmdRouter.POST("/record/log", baseApi.LoadRecordLog)
		cmdRouter.POST("/record/clean", baseApi.CleanRecord)
	}
}
