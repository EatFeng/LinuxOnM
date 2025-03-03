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

		containerRouter.POST("", baseApi.ContainerCreate)
		containerRouter.POST("/update", baseApi.ContainerUpdate)
		containerRouter.POST("/upgrade", baseApi.ContainerUpgrade)
		containerRouter.POST("/info", baseApi.ContainerInfo)
		containerRouter.POST("/search", baseApi.SearchContainer)
		containerRouter.POST("/list", baseApi.ListContainer)
		containerRouter.GET("/list/stats", baseApi.ContainerListStats)
		containerRouter.GET("/search/log", baseApi.ContainerLogs)
		containerRouter.POST("/download/log", baseApi.DownloadContainerLogs)
		containerRouter.POST("/clean/log", baseApi.CleanContainerLog)
		containerRouter.GET("/limit", baseApi.LoadResourceLimit)
		containerRouter.POST("/operate", baseApi.ContainerOperation)
		containerRouter.POST("/inspect", baseApi.Inspect)
		containerRouter.POST("/rename", baseApi.ContainerRename)
		containerRouter.POST("/commit", baseApi.ContainerCommit)
		containerRouter.POST("/prune", baseApi.ContainerPrune)

		containerRouter.GET("/repo", baseApi.ListRepo)

		containerRouter.GET("/image", baseApi.ListImage)
		containerRouter.GET("/image/all", baseApi.ListAllImage)
		containerRouter.POST("/image/search", baseApi.SearchImage)

		containerRouter.GET("/network", baseApi.ListNetwork)

		containerRouter.GET("/volume", baseApi.ListVolume)

		containerRouter.GET("/docker/status", baseApi.LoadDockerStatus)
	}
}
