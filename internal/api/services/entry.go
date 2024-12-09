package services

import "LinuxOnM/internal/repositories"

var (
	logRepo     = repositories.NewLogRepository()
	commonRepo  = repositories.NewCommonRepository()
	settingRepo = repositories.NewISettingRepo()
	hostRepo    = repositories.NewIHostRepo()
	groupRepo   = repositories.NewIGroupRepo()
	commandRepo = repositories.NewICommandRepo()
	cronjobRepo = repositories.NewICronjobRepo()
)
