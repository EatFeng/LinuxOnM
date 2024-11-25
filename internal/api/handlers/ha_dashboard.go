package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"errors"
	"github.com/gin-gonic/gin"
)

// @Tags Dashboard
// @Summary Load os info
// @Description 获取服务器基础数据
// @Accept json
// @Success 200 {object} dto.OsInfo
// @Security ApiKeyAuth
// @Router /dashboard/base/os [get]
func (b *BaseApi) LoadDashboardOsInfo(c *gin.Context) {
	data, err := dashboardService.LoadOsInfo()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, data)
}

// @Tags Dashboard
// @Summary Load dashboard current info
// @Description 获取首页实时数据
// @Accept json
// @Param ioOption path string true "request"
// @Param netOption path string true "request"
// @Success 200 {object} dto.DashboardCurrent
// @Security ApiKeyAuth
// @Router /dashboard/current/:ioOption/:netOption [get]
func (b *BaseApi) LoadDashboardCurrentInfo(c *gin.Context) {
	ioOption, ok := c.Params.Get("ioOption")
	if !ok {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, errors.New("error ioOption in path"))
		return
	}
	netOption, ok := c.Params.Get("netOption")
	if !ok {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, errors.New("error netOption in path"))
		return
	}

	data := dashboardService.LoadCurrentInfo(ioOption, netOption)
	helper.SuccessWithData(c, data)
}
