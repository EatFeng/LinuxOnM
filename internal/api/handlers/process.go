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
// @Summary Stop Process
// @Param request body request.ProcessReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /process/stop [post]
// @x-panel-log {"bodyKeys":["PID"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"结束进程 [PID]","formatEN":"结束进程 [PID]"}
func (b *BaseApi) StopProcess(c *gin.Context) {
	var req request.ProcessReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := processService.StopProcess(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	helper.SuccessWithOutData(c)
}
