package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"github.com/gin-gonic/gin"
)

// LoadDockerStatus
// @Tags Container Docker
// @Summary Load docker status
// @Description 获取 docker 服务状态
// @Produce json
// @Success 200 {string} status
// @Security ApiKeyAuth
// @Router /container/docker/status [get]
func (b *BaseApi) LoadDockerStatus(c *gin.Context) {
	status := dockerService.LoadDockerStatus()
	helper.SuccessWithData(c, status)
}
