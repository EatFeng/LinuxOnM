package cron

import (
	"LinuxOnM/internal/api/services"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/common"
	"LinuxOnM/internal/utils/ntp"
	"time"

	"github.com/robfig/cron/v3"
)

func Run() {
	nyc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	global.Cron = cron.New(cron.WithLocation(nyc), cron.WithChain(cron.Recover(cron.DefaultLogger)), cron.WithChain(cron.DelayIfStillRunning(cron.DefaultLogger)))

	var (
		interval models.Setting
		status   models.Setting
	)
	go syncBeforeStart()
	if err := global.DB.Where("key = ?", "MonitorStatus").Find(&status).Error; err != nil {
		global.LOG.Errorf("load monitor status from db failed, err: %v", err)
	}
	if status.Value == "enable" {
		if err := global.DB.Where("key = ?", "MonitorInterval").Find(&interval).Error; err != nil {
			global.LOG.Errorf("load monitor interval from db failed, err: %v", err)
		}
		if err := services.StartMonitor(false, interval.Value); err != nil {
			global.LOG.Errorf("can not add monitor corn job: %s", err.Error())
		}
	}

	global.Cron.Start()
}

func syncBeforeStart() {
	var ntpSite models.Setting
	if err := global.DB.Where("key = ?", "NtpSite").Find(&ntpSite).Error; err != nil {
		global.LOG.Errorf("load ntp serve from db failed, err: %v", err)
	}
	if len(ntpSite.Value) == 0 {
		ntpSite.Value = "pool.ntp.org"
	}
	ntime, err := ntp.GetRemoteTime(ntpSite.Value)
	if err != nil {
		global.LOG.Errorf("load remote time with [%s] failed, err: %v", ntpSite.Value, err)
		return
	}
	ts := ntime.Format(constant.DateTimeLayout)
	if err := ntp.UpdateSystemTime(ts); err != nil {
		global.LOG.Errorf("failed to synchronize system time with [%s], err: %v", ntpSite.Value, err)
	}
	global.LOG.Debugf("synchronize system time with [%s] successful!", ntpSite.Value)
}
