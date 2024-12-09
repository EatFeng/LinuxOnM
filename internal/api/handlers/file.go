package handlers

import (
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

// @Tags File
// @Summary Read file by Line
// @Description 按行读取日志文件
// @Param request body request.FileReadByLineReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/read [post]
func (b *BaseApi) ReadFileByLine(c *gin.Context) {
	var req request.FileReadByLineReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	res, err := fileService.ReadLogByLine(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, res)
}

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
