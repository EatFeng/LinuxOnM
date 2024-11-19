package routers

import (
	"github.com/gin-gonic/gin"
	"linux-ops-platform/internal/api/handlers"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 定义API版本
	api := r.Group("/api/v1")

	// 定义路由
	api.GET("/home", handlers.HomeHandler)
}
