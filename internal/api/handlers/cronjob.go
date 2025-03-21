package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/common"
	"github.com/gin-gonic/gin"
	"time"
)

// CreateCronjob
// @Tags Cronjob
// @Summary Create cronjob
// @Description 创建计划任务
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
// @Description 获取计划任务分页
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
// @Description 更新计划任务
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

// UpdateCronjobStatus
// @Tags Cronjob
// @Summary Update cronjob status
// @Description 更新计划任务状态
// @Accept json
// @Param request body dto.CronjobUpdateStatus true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob/status [post]
// @x-panel-log {"bodyKeys":["id","status"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"cronjobs","output_column":"name","output_value":"name"}],"formatZH":"修改计划任务 [name] 状态为 [status]","formatEN":"change the status of cronjob [name] to [status]."}
func (b *BaseApi) UpdateCronjobStatus(c *gin.Context) {
	var req dto.CronjobUpdateStatus
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := cronjobService.UpdateStatus(req.ID, req.Status); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// HandleOnce
// @Tags Cronjob
// @Summary Handle cronjob once
// @Description 手动执行计划任务
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob/handle [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"cronjobs","output_column":"name","output_value":"name"}],"formatZH":"手动执行计划任务 [name]","formatEN":"manually execute the cronjob [name]"}
func (b *BaseApi) HandleOnce(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := cronjobService.HandleOnce(req.ID); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// SearchJobRecords
// @Tags Cronjob
// @Summary Page job records
// @Description 获取计划任务记录
// @Accept json
// @Param request body dto.SearchRecord true "request"
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /cronjob/record/search [post]
func (b *BaseApi) SearchJobRecords(c *gin.Context) {
	var req dto.SearchRecord
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	loc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	req.StartTime = req.StartTime.In(loc)
	req.EndTime = req.EndTime.In(loc)

	total, list, err := cronjobService.SearchRecords(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// LoadRecordLog
// @Tags Cronjob
// @Summary Load Cronjob record log
// @Description 获取计划任务记录日志
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob/record/log [post]
func (b *BaseApi) LoadRecordLog(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	content := cronjobService.LoadRecordLog(req)
	helper.SuccessWithData(c, content)
}

// CleanRecord
// @Tags Cronjob
// @Summary Clean job records
// @Description 清空计划任务记录
// @Accept json
// @Param request body dto.CronjobClean true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob/record/clean [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"cronjobs","output_column":"name","output_value":"name"}],"formatZH":"清空计划任务记录 [name]","formatEN":"clean cronjob [name] records"}
func (b *BaseApi) CleanRecord(c *gin.Context) {
	var req dto.CronjobClean
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := cronjobService.CleanRecord(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, nil)
}

// DeleteCronjob
// @Tags Cronjob
// @Summary Delete cronjob
// @Description 删除计划任务
// @Accept json
// @Param request body dto.CronjobBatchDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /cronjob/del [post]
// @x-panel-log {"bodyKeys":["ids"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"ids","isList":true,"db":"cronjobs","output_column":"name","output_value":"names"}],"formatZH":"删除计划任务 [names]","formatEN":"delete cronjob [names]"}
func (b *BaseApi) DeleteCronjob(c *gin.Context) {
	var req dto.CronjobBatchDelete
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := cronjobService.Delete(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
