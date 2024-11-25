package helper

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
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
		case errors.As(&buserr.BusinessError{}, err):
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
	return nil
}
