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
		containerRouter.GET("/exec", baseApi.ContainerWsSsh)
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
		containerRouter.POST("/load/log", baseApi.LoadContainerLog)
		containerRouter.GET("/limit", baseApi.LoadResourceLimit)
		containerRouter.POST("/operate", baseApi.ContainerOperation)
		containerRouter.POST("/inspect", baseApi.Inspect)
		containerRouter.POST("/rename", baseApi.ContainerRename)
		containerRouter.POST("/commit", baseApi.ContainerCommit)
		containerRouter.POST("/prune", baseApi.ContainerPrune)

		containerRouter.GET("/repo", baseApi.ListRepo)
		containerRouter.POST("/repo/status", baseApi.CheckRepoStatus)
		containerRouter.POST("/repo/search", baseApi.SearchRepo)
		containerRouter.POST("/repo", baseApi.CreateRepo)
		containerRouter.POST("/repo/update", baseApi.UpdateRepo)
		containerRouter.POST("/repo/del", baseApi.DeleteRepo)

		containerRouter.GET("/image", baseApi.ListImage)
		containerRouter.GET("/image/all", baseApi.ListAllImage)
		containerRouter.POST("/image/search", baseApi.SearchImage)
		containerRouter.POST("/image/pull", baseApi.ImagePull)
		containerRouter.POST("/image/push", baseApi.ImagePush)
		containerRouter.POST("/image/save", baseApi.ImageSave)
		containerRouter.POST("/image/load", baseApi.ImageLoad)
		containerRouter.POST("/image/remove", baseApi.ImageRemove)
		containerRouter.POST("/image/tag", baseApi.ImageTag)
		containerRouter.POST("/image/build", baseApi.ImageBuild)

		containerRouter.GET("/network", baseApi.ListNetwork)
		containerRouter.POST("/network/del", baseApi.DeleteNetwork)
		containerRouter.POST("/network/search", baseApi.SearchNetwork)
		containerRouter.POST("/network", baseApi.CreateNetwork)

		containerRouter.GET("/volume", baseApi.ListVolume)
		containerRouter.POST("/volume/del", baseApi.DeleteVolume)
		containerRouter.POST("/volume/search", baseApi.SearchVolume)
		containerRouter.POST("/volume", baseApi.CreateVolume)

		containerRouter.GET("/template", baseApi.ListComposeTemplate)
		containerRouter.POST("/template", baseApi.CreateComposeTemplate)
		containerRouter.POST("/template/search", baseApi.SearchComposeTemplate)
		containerRouter.POST("/template/update", baseApi.UpdateComposeTemplate)
		containerRouter.POST("/template/del", baseApi.DeleteComposeTemplate)

		containerRouter.POST("/compose/search", baseApi.SearchCompose)
		containerRouter.POST("/compose", baseApi.CreateCompose)
		containerRouter.POST("/compose/test", baseApi.TestCompose)
		containerRouter.POST("/compose/operate", baseApi.OperatorCompose)
		containerRouter.POST("/compose/update", baseApi.ComposeUpdate)
		containerRouter.GET("/compose/search/log", baseApi.ComposeLogs)

		containerRouter.GET("/docker/status", baseApi.LoadDockerStatus)
		containerRouter.POST("/docker/operate", baseApi.OperateDocker)
		containerRouter.GET("/daemonjson", baseApi.LoadDaemonJson)
		containerRouter.GET("/daemonjson/file", baseApi.LoadDaemonJsonFile)
		containerRouter.POST("/daemonjson/update", baseApi.UpdateDaemonJson)
		containerRouter.POST("/logoption/update", baseApi.UpdateLogOption)
		containerRouter.POST("/ipv6option/update", baseApi.UpdateIpv6Option)
		containerRouter.POST("/daemonjson/update/byfile", baseApi.UpdateDaemonJsonByFile)
	}
}
