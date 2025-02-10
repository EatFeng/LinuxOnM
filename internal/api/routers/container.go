package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"
	"github.com/gin-gonic/gin"
)

type ContainerRouter struct{}

func (s *ContainerRouter) InitRouter(Router *gin.RouterGroup) {
	containerRouter := Router.Group("container").Use(middleware.PasswordExpired())
	baseApi := handler.ApiGroupApp.BaseApi
	{
		containerRouter.GET("/stats/:id", baseApi.ContainerStats)

		containerRouter.POST("/info", baseApi.ContainerInfo)
		containerRouter.POST("/search", baseApi.SearchContainer)
		containerRouter.POST("/list", baseApi.ListContainer)
		containerRouter.GET("/list/stats", baseApi.ContainerListStats)
		containerRouter.GET("/search/log", baseApi.ContainerLogs)
		containerRouter.GET("/limit", baseApi.LoadResourceLimit)

		containerRouter.GET("/repo", baseApi.ListRepo)

		containerRouter.GET("/image", baseApi.ListImage)

		containerRouter.GET("/network", baseApi.ListNetwork)

		containerRouter.GET("/volume", baseApi.ListVolume)

		containerRouter.GET("/docker/status", baseApi.LoadDockerStatus)
	}
}
