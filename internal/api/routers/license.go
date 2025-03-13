package routers

import (
	"LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"

	"github.com/gin-gonic/gin"
)

type LicenseRouter struct{}

func (s *LicenseRouter) InitRouter(Router *gin.RouterGroup) {
	licenseRouter := Router.Group("license").Use(middleware.PasswordExpired())
	baseApi := handlers.ApiGroupApp.BaseApi
	{
		licenseRouter.POST("/upload", baseApi.HandleLicenseUpload)
		licenseRouter.GET("/status", baseApi.GetLicenseInfo)
	}
}
