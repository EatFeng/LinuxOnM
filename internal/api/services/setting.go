package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/encrypt"
	"encoding/json"
	"time"
)

type SettingService struct{}

type ISettingService interface {
	GetSettingInfo() (*dto.SettingInfo, error)
}

func NewISettingService() ISettingService {
	return &SettingService{}
}

func (u *SettingService) GetSettingInfo() (*dto.SettingInfo, error) {
	settings, err := settingRepo.GetList()
	if err != nil {
		return nil, constant.ErrRecordNotFound
	}
	settingMap := make(map[string]string)
	for _, setting := range settings {
		settingMap[setting.Key] = setting.Value
	}
	var info dto.SettingInfo
	arr, err := json.Marshal(settingMap)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(arr, &info); err != nil {
		return nil, err
	}
	if info.ProxyPasswdKeep != constant.StatusEnable {
		info.ProxyPasswd = ""
	} else {
		info.ProxyPasswd, _ = encrypt.StringDecrypt(info.ProxyPasswd)
	}

	info.LocalTime = time.Now().Format("2006-01-02 15:04:05 MST -0700")
	return &info, err
}
