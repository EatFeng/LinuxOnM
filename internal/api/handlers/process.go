package handlers

import (
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	websocket2 "LinuxOnM/internal/utils/websocket"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) ProcessWs(c *gin.Context) {
	ws, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	wsClient := websocket2.NewWsClient("processClient", ws)
	go wsClient.Read()
	go wsClient.Write()
}

// @Tags Process
// @Summary Kill Process
// @Param request body request.ProcessReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /process/kill [post]
// @x-panel-log {"bodyKeys":["PID"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"结束进程 [PID]","formatEN":"结束进程 [PID]"}
func (b *BaseApi) KillProcess(c *gin.Context) {
	var req request.ProcessReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := processService.KillProcess(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)
}

func (b *BaseApi) GetProcessContent(c *gin.Context) {
	var req request.ProcessRequest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	info, err := processService.GetProcessContent(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithData(c, info)
}

func (b *BaseApi) StartProcess(c *gin.Context) {
	var req request.ProcessRequest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := processService.StartProcess(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)
}

func (b *BaseApi) StopProcess(c *gin.Context) {
	var req request.ProcessRequest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := processService.StopProcess(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)
}

func (b *BaseApi) EnableProcess(c *gin.Context) {
	var req request.ProcessRequest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := processService.EnableProcess(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)
}

func (b *BaseApi) DisableProcess(c *gin.Context) {
	var req request.ProcessRequest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := processService.DisableProcess(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)
}

func (b *BaseApi) StatusProcess(c *gin.Context) {
	var req request.ProcessRequest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	status, err := processService.StatusProcess(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithData(c, status)
}

func (b *BaseApi) CreateProcess(c *gin.Context) {
	var req request.ProcessCreate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	err := processService.CreateProcess(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)

}
