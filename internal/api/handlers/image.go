package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"

	"github.com/gin-gonic/gin"
)

// ListImage
// @Tags Container Image
// @Summary load images options
// @Description 获取镜像名称列表
// @Produce json
// @Success 200 {array} dto.Options
// @Security ApiKeyAuth
// @Router /container/image [get]
func (b *BaseApi) ListImage(c *gin.Context) {
	list, err := imageService.List()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, list)
}

// @Tags Container Image
// @Summary List all images
// @Description 获取所有镜像列表
// @Produce json
// @Success 200 {array} dto.ImageInfo
// @Security ApiKeyAuth
// @Router /containers/image/all [get]
func (b *BaseApi) ListAllImage(c *gin.Context) {
	list, err := imageService.ListAll()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, list)
}
