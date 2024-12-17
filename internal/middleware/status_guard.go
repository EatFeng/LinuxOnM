package middleware

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/repositories"
	"github.com/gin-gonic/gin"
)

func StatusGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		settingRepo := repositories.NewISettingRepo()
		status, err := settingRepo.Get(settingRepo.WithByKey("SystemStatus"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}
		if status.Value != "Free" {
			helper.ErrorWithDetail(c, constant.CodeGlobalLoading, status.Value, err)
			return
		}
		c.Next()
	}
}
