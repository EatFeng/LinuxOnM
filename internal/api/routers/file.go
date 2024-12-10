package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type FileRouter struct{}

func (f *FileRouter) InitRouter(Router *gin.RouterGroup) {
	fileRouter := Router.Group("file")
	baseApi := handler.ApiGroupApp.BaseApi
	{
		fileRouter.POST("/read", baseApi.ReadFileByLine)
		fileRouter.POST("/search", baseApi.ListFiles)
	}
}
