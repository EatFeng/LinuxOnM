package middleware

import (
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/repositories"
	"errors"
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

func WhiteAllow() gin.HandlerFunc {
	return func(c *gin.Context) {
		settingRepo := repositories.NewISettingRepo()
		status, err := settingRepo.Get(settingRepo.WithByKey("AllowIPs"))
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}

		if len(status.Value) == 0 {
			c.Next()
			return
		}
		clientIP := c.ClientIP()
		for _, ip := range strings.Split(status.Value, ",") {
			if len(ip) == 0 {
				continue
			}
			if ip == clientIP || (strings.Contains(ip, "/") && checkIpInCidr(ip, clientIP)) {
				c.Next()
				return
			}
		}
		if LoadErrCode("err-ip") != 200 {
			helper.ErrResponse(c, LoadErrCode("err-ip"))
			return
		}
		helper.ErrorWithDetail(c, constant.CodeErrIP, constant.ErrTypeInternalServer, errors.New("IP address not allowed"))
	}
}

func checkIpInCidr(cidr, checkIP string) bool {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		global.LOG.Errorf("parse CIDR %s failed, err: %v", cidr, err)
		return false
	}
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		if ip.String() == checkIP {
			return true
		}
	}
	return false
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
