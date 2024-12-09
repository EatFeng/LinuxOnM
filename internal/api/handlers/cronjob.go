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

// SearchCronjob
// @Tags Cronjob
// @Summary Page cronjob
// @Description Get page of cronjob
// @Accept json
// @Param request body dto.PageCronjob true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /cronjob/search [post]
func (b *BaseApi) SearchCronjob(c *gin.Context) {
	var req dto.PageCronjob
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	total, list, err := cronjobService.SearchWithPage(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// UpdateCronjob
// @Tags Cronjob
// @Summary Update cronjob
// @Description This function is used to update an existing cronjob. It first binds and validates the incoming JSON request body of type dto.CronjobUpdate. Then, it attempts to update the corresponding cronjob record in the system.
// @Accept json
// @Param request body dto.CronjobUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob/update [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"cronjobs","output_column":"name","output_value":"name"}],"formatZH":"更新计划任务 [name]","formatEN":"update cronjob [name]"}
func (b *BaseApi) UpdateCronjob(c *gin.Context) {
	var req dto.CronjobUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := cronjobService.Update(req.ID, req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
