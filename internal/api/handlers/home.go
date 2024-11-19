package handlers

import (
	"github.com/gin-gonic/gin"
	"linux-ops-platform/internal/services"
	"net/http"
)

// HomeHandler 处理首页请求
func HomeHandler(c *gin.Context) {
	// 调用服务层获取数据
	data, err := services.GetHomeData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回响应
	c.JSON(http.StatusOK, data)
}
