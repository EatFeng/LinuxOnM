package middleware

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/repositories"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

func BindDomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		settingRepo := repositories.NewISettingRepo()
		status, err := settingRepo.Get(settingRepo.WithByKey("BindDomain"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}
		if len(status.Value) == 0 {
			c.Next()
			return
		}
		domains := c.Request.Host
		parts := strings.Split(c.Request.Host, ":")
		if len(parts) > 0 {
			domains = parts[0]
		}

		if domains != status.Value {
			helper.ErrorWithDetail(c, constant.CodeErrDomain, constant.ErrTypeInternalServer, errors.New("domain not allowed"))
			return
		}
		c.Next()
	}
}
