package handlers

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/api/services"
	"LinuxOnM/internal/constant"

	"github.com/gin-gonic/gin"
)

func (b *BaseApi) GetCurrentStatus(c *gin.Context) {
	service := services.NewSystemStatusService()
	status, err := service.GetCurrentStatus()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, status)
}
