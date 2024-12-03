package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// @Tags Logs
// @Summary Page login logs
// @Description 获取系统登录日志列表分页
// @Accept json
// @Param request body dto.SearchLoginLogWithPage true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /logs/login [post]
func (b *BaseApi) GetLoginLog(c *gin.Context) {
	var req dto.SearchLoginLogWithPage
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	total, list, err := logService.PageLoginLog(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// @Tags Logs
// @Summary Page operation logs
// @Description 获取系统操作日志列表分页
// @Accept json
// @Param request body dto.SearchOpLogWithPage true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /log/operation [post]
func (b *BaseApi) GetOperationLog(c *gin.Context) {
	var req dto.SearchOpLogWithPage
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	total, list, err := logService.PageOperationLog(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// @Tags Logs
// @Summary Load system log files
// @Description 获取系统日志文件列表
// @Success 200
// @Security ApiKeyAuth
// @Router /log/system/files [get]
func (b *BaseApi) GetSystemFiles(c *gin.Context) {
	data, err := logService.ListSystemLogFile()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, data)
}

// @Tags SSH
// @Summary Load host SSH logs
// @Description 获取 SSH 登录日志
// @Accept json
// @Param request body dto.SearchSSHLog true "request"
// @Success 200 {object} dto.SSHLog
// @Security ApiKeyAuth
// @Router /log/ssh [post]
func (b *BaseApi) LoadSSHLog(c *gin.Context) {
	var req dto.SearchSSHLog
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	data, err := logService.LoadSSHLog(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, data)
}
