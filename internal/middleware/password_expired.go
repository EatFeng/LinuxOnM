package middleware

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/repositories"
	"LinuxOnM/internal/utils/common"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func PasswordExpired() gin.HandlerFunc {
	return func(c *gin.Context) {
		settingRepo := repositories.NewISettingRepo()
		setting, err := settingRepo.Get(settingRepo.WithByKey("ExpirationDays"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypePasswordExpired, err)
			return
		}
		expiredDays, _ := strconv.Atoi(setting.Value)
		if expiredDays == 0 {
			c.Next()
			return
		}

		extime, err := settingRepo.Get(settingRepo.WithByKey("ExpirationTime"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypePasswordExpired, err)
			return
		}
		loc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
		expiredTime, err := time.ParseInLocation(constant.DateTimeLayout, extime.Value, loc)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodePasswordExpired, constant.ErrTypePasswordExpired, err)
			return
		}
		if time.Now().After(expiredTime) {
			helper.ErrorWithDetail(c, constant.CodePasswordExpired, constant.ErrTypePasswordExpired, err)
			return
		}
		c.Next()
	}
}
