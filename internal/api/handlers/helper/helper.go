package helper

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

func ErrorWithDetail(ctx *gin.Context, code int, msgKey string, err error) {
	res := dto.Response{
		Code:    code,
		Message: "",
	}
	if msgKey == constant.ErrTypeInternalServer {
		switch {
		case errors.Is(constant.ErrRecordExist, err):
			res.Message = "Record already exists"
		case errors.Is(constant.ErrRecordNotFound, err):
			res.Message = "Record not found"
		case errors.Is(constant.ErrInvalidParams, err):
			res.Message = "Invalid parameters"
		case errors.Is(constant.ErrStructTransform, err):
			res.Message = fmt.Sprintf("Struct transform error: %v", err)
		case errors.Is(constant.ErrCaptchaCode, err):
			res.Code = constant.CodeAuth
			res.Message = "Captcha code error"
		case errors.Is(constant.ErrAuth, err):
			res.Code = constant.CodeAuth
			res.Message = "Auth error"
		case errors.Is(constant.ErrInitialPassword, err):
			res.Message = "Initial Password error"
		case errors.As(err, &buserr.BusinessError{}):
			res.Message = err.Error()
		default:
			res.Message = fmt.Sprintf("%s: %v", msgKey, err)
		}
	} else {
		res.Message = fmt.Sprintf("%s: %v", msgKey, err)
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func SuccessWithData(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := dto.Response{
		Code: constant.CodeSuccess,
		Data: data,
	}
	ctx.JSON(http.StatusOK, res)
	ctx.Abort()
}

func CheckBindAndValidate(ctx *gin.Context, req interface{}) error {
	if err := ctx.ShouldBindJSON(req); err != nil {
		ErrorWithDetail(ctx, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return err
	}
	if err := global.VALID.Struct(req); err != nil {
		ErrorWithDetail(ctx, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return err
	}
	return nil
}

func GetParamID(c *gin.Context) (uint, error) {
	idParam, ok := c.Params.Get("id")
	if !ok {
		return 0, errors.New("error id in path")
	}
	intNum, _ := strconv.Atoi(idParam)
	return uint(intNum), nil
}
