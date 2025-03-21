package middleware

import (
	"LinuxOnM/internal/repositories"
	"net/http"
)

func LoadErrCode(errInfo string) int {
	settingRepo := repositories.NewISettingRepo()
	codeVal, err := settingRepo.Get(settingRepo.WithByKey("NoAuthSetting"))
	if err != nil {
		return 500
	}

	switch codeVal.Value {
	case "400":
		return http.StatusBadRequest
	case "401":
		return http.StatusUnauthorized
	case "403":
		return http.StatusForbidden
	case "404":
		return http.StatusNotFound
	case "408":
		return http.StatusRequestTimeout
	case "416":
		return http.StatusRequestedRangeNotSatisfiable
	default:
		return http.StatusOK
	}
}
