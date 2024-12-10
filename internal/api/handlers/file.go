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
// @Summary List files
// @Accept json
// @Param request body request.FileOption true "request"
// @Success 200 {object} response.FileInfo
// @Security ApiKeyAuth
// @Router /files/search [post]
func (b *BaseApi) ListFiles(c *gin.Context) {
	var req request.FileOption
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	files, err := fileService.GetFileList(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, files)
}

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
