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

// LoadInterfaceAddr
// @Tags System Setting
// @Summary Load system address
// @Description 获取系统IPv4和IPv6地址信息
// @Accept json
// @Success 200
// @Security ApiKeyAuth
// @Router /setting/interface [get]
func (b *BaseApi) LoadInterfaceAddr(c *gin.Context) {
	data, err := settingService.LoadInterfaceAddr()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, data)
}

// GetSystemAvailable
// @Tags System Setting
// @Summary Load system available status
// @Description 获取系统可用状态
// @Success 200
// @Security ApiKeyAuth
// @Router /setting/search/available [get]
func (b *BaseApi) GetSystemAvailable(c *gin.Context) {
	helper.SuccessWithData(c, nil)
}

// UpdatePassword
// @Tags System Setting
// @Summary Update system password
// @Description 更新系统登录密码
// @Accept json
// @Param request body dto.PasswordUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /setting/update/password [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFunctions":[],"formatZH":"修改系统密码","formatEN":"update system password"}
func (b *BaseApi) UpdatePassword(c *gin.Context) {
	var req dto.PasswordUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := settingService.UpdatePassword(c, req.OldPassword, req.NewPassword); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// HandlePasswordExpired
// @Tags System Setting
// @Summary Reset system password expired
// @Description 重置过期系统登录密码
// @Accept json
// @Param request body dto.PasswordUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /setting/expired/handle [post]
// @x-panel-log {"bodyKeys":[],"paramKeys":[],"BeforeFunctions":[],"formatZH":"重置过期密码","formatEN":"reset an expired Password"}
func (b *BaseApi) HandlePasswordExpired(c *gin.Context) {
	var req dto.PasswordUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := settingService.HandlePasswordExpired(c, req.OldPassword, req.NewPassword); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
