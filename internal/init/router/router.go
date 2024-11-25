package router

import (
	rou "LinuxOnM/internal/api/routers"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routers() *gin.Engine {
	Router := gin.Default()

	Router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "The requested URL was not found on this server.",
		})
	})

	PublicGroup := Router.Group("")
	{
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		PublicGroup.Use(gzip.Gzip(gzip.DefaultCompression))
	}

	PrivateGroup := Router.Group("/api/handler")
	for _, router := range rou.RouterGroupApp {
		router.InitRouter(PrivateGroup)
	}

	return Router
}
