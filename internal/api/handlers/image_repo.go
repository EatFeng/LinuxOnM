package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// ListRepo
// @Tags Container Image-repo
// @Summary List image repos
// @Description 获取镜像仓库列表
// @Produce json
// @Success 200 {array} dto.ImageRepoOption
// @Security ApiKeyAuth
// @Router /container/repo [get]
func (b *BaseApi) ListRepo(c *gin.Context) {
	list, err := imageRepoService.List()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, list)
}
