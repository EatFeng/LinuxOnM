package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/middleware"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/captcha"
	"LinuxOnM/internal/utils/qqwry"
	"encoding/base64"
	"github.com/gin-gonic/gin"
)

// @Tags Auth
// @Summary User login
// @Description 用户登录
// @Accept json
// @Param EntranceCode header string true "安全入口 base64 加密串"
// @Param request body dto.Login true "request"
// @Success 200 {object} dto.UserLoginInfo
// @Router /auth/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var req dto.Login
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if req.AuthMethod != "jwt" && !req.IgnoreCaptcha {
		if err := captcha.VerifyCode(req.CaptchaID, req.Captcha); err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}
	}
	entranceItem := c.Request.Header.Get("entranceCode")
	var entrance []byte
	if len(entranceItem) != 0 {
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

// CheckIsSafety
// @Tags Auth
// @Summary Load safety status
// @Description 获取系统安全登录状态
// @Success 200
// @Router /auth/is-safety [get]
func (b *BaseApi) CheckIsSafety(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	status, err := authService.CheckIsSafety(code)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	if status == "disable" && len(code) != 0 {
		helper.ErrorWithDetail(c, constant.CodeErrNotFound, constant.ErrTypeInternalServer, err)
		return
	}
	if status == "unpass" {
		if middleware.LoadErrCode("err-entrance") != 200 {
			helper.ErrResponse(c, middleware.LoadErrCode("err-entrance"))
			return
		}
		helper.ErrorWithDetail(c, constant.CodeErrEntrance, constant.ErrTypeInternalServer, nil)
		return
	}
	helper.SuccessWithOutData(c)
}

func saveLoginLog(c *gin.Context, err error) {
	var logs models.LoginLog
	if err != nil {
		logs.Status = constant.StatusFailed
		logs.Message = err.Error()
	} else {
		logs.Status = constant.StatusSuccess
	}
	logs.IP = c.ClientIP()
	qqWry, err := qqwry.NewQQwry()
	if err != nil {
		global.LOG.Errorf("load qqwry datas failed: %s", err)
	}
	res := qqWry.Find(logs.IP)
	logs.Agent = c.GetHeader("User-Agent")
	logs.Address = res.Area
	_ = logService.CreateLoginLog(logs)
}
