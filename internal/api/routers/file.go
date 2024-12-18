package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"
	"github.com/gin-gonic/gin"
)

type FileRouter struct{}

func (f *FileRouter) InitRouter(Router *gin.RouterGroup) {
	fileRouter := Router.Group("file").Use(middleware.PasswordExpired())
	baseApi := handler.ApiGroupApp.BaseApi
	{
		fileRouter.POST("/read", baseApi.ReadFileByLine)
		fileRouter.POST("/search", baseApi.ListFiles)
		fileRouter.POST("", baseApi.CreateFile)
		fileRouter.POST("/del", baseApi.DeleteFile)
		fileRouter.POST("/upload", baseApi.UploadFiles)
		fileRouter.POST("/content", baseApi.GetContent)
		fileRouter.POST("/tree", baseApi.GetFileTree)
		fileRouter.POST("/batch/role", baseApi.BatchChangeModeAndOwner)
		fileRouter.POST("/check", baseApi.CheckFile)
		fileRouter.POST("/move", baseApi.MoveFile)
		fileRouter.POST("/rename", baseApi.ChangeFileName)
		fileRouter.POST("/compress", baseApi.CompressFile)
	}
}
