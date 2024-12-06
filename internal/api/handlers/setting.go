package handlers

import (
	"LinuxOnM/internal/api/dto"
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

// UpdateSetting
// @Tags System Setting
// @Summary Update system settings.
// @Description Modify the system settings.
// @Accept json
// @Param request body dto.SettingUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /setting/update [post]
// @x-panel-log {"bodyKeys":["key","value"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"修改系统配置 [key] => [value]","formatEN":"update system setting [key] => [value]"}
func (b *BaseApi) UpdateSetting(c *gin.Context) {
	var req dto.SettingUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := settingService.Update(req.Key, req.Value); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
