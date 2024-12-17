package router

import (
	"LinuxOnM/docs"
	rou "LinuxOnM/internal/api/routers"
	"LinuxOnM/internal/middleware"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers() *gin.Engine {
	Router := gin.Default()
	Router.Use(middleware.OperationLog())

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The requested URL was not found on this server.",
		})
	})

	swaggerRouter := Router.Group("linuxonm")
	docs.SwaggerInfo.BasePath = "/api/handler"
	swaggerRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	PublicGroup := Router.Group("")
	{
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		PublicGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	PrivateGroup := Router.Group("/api/handler")
	PrivateGroup.Use(middleware.WhiteAllow())
	PrivateGroup.Use(middleware.BindDomain())
	PrivateGroup.Use(middleware.StatusGuard())
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	return Router
}
