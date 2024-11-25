package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"encoding/base64"
	"github.com/gin-gonic/gin"
)

func (b *BaseApi) Login(c *gin.Context) {
	var req dto.Login
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if req.AuthMethod != "jwt" && !req.IgnoreCapycha {
		if err := captcha.VerifyCode(req.CaptchaID, req.Captcha); err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}
	}
	entranceItem := c.Request.Header.Get("entranceCode")
	var entrance []byte
	if len(entranceItem) == 0 {
		entrance, _ = base64.StdEncoding.DecodeString(entranceItem)
	}

	user, err := authService.Login(c, req, string(entrance))
	go saveLoginLog(c, err)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, user)
}
