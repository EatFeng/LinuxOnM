package hook

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/repositories"
)

func Init() {
	settingRepo := repositories.NewISettingRepo()
	portSetting, err := settingRepo.Get(settingRepo.WithByKey("ServerPort"))
	if err != nil {
		global.LOG.Errorf("load service port from setting failed, err: %v", err)
	}
	global.CONF.System.Port = portSetting.Value

	ipv6Setting, err := settingRepo.Get(settingRepo.WithByKey("Ipv6"))
	if err != nil {
		global.LOG.Errorf("load ipv6 status from setting failed, err: %v", err)
	}
	global.CONF.System.Ipv6 = ipv6Setting.Value

	bindAddressSetting, err := settingRepo.Get(settingRepo.WithByKey("BindAddress"))
	if err != nil {
		global.LOG.Errorf("load bind address from setting failed, err: %v", err)
	}
	global.CONF.System.BindAddress = bindAddressSetting.Value

	sslSetting, err := settingRepo.Get(settingRepo.WithByKey("SSL"))
	if err != nil {
		global.LOG.Errorf("load service ssl from setting failed, err: %v", err)
	}
	global.CONF.System.SSL = sslSetting.Value

	if _, err := settingRepo.Get(settingRepo.WithByKey("SystemStatus")); err != nil {
		_ = settingRepo.Create("SystemStatus", "Free")
	}
	if err := settingRepo.Update("SystemStatus", "Free"); err != nil {
		global.LOG.Fatalf("init service before start failed, err: %v", err)
	}

	handleCronjobStatus()
}

func handleCronjobStatus() {
	_ = global.DB.Model(&models.JobRecords{}).Where("status = ?", constant.StatusWaiting).
		Updates(map[string]interface{}{
			"status":  constant.StatusFailed,
			"message": "the task was interrupted due to the restart of the myapp_LinuxOnM service",
		}).Error
}
