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
// @Param request body dto.SearchLgLogWithPage true "request"
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
