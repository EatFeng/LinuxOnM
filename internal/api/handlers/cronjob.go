package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// CreateCronjob
// @Tags Cronjob
// @Summary Create cronjob
// @Description Create a Cronjob
// @Accept json
// @Param request body dto.CronjobCreate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob [post]
// @x-panel-log {"bodyKeys":["type","name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建计划任务 [type][name]","formatEN":"create cronjob [type][name]"}
func (b *BaseApi) CreateCronjob(c *gin.Context) {
	var req dto.CronjobCreate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := cronjobService.Create(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
