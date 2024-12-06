package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// GetSettingInfo
// @Tags System Setting
// @Summary Load system setting information.
// @Description Retrieve the system setting information.
// @Success 200 {object} dto.SettingInfo
// @Security ApiKeyAuth
// @Router /setting/search [post]
func (b *BaseApi) GetSettingInfo(c *gin.Context) {
	setting, err := settingService.GetSettingInfo()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, setting)
}
