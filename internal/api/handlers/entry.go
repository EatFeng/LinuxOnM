package handlers

import "LinuxOnM/internal/api/services"

type BaseApi struct{}

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

var (
	dashboardService = services.NewDashboardService()
	logService       = services.NewILogService()
	authService      = services.NewIAuthService()
	hostService      = services.NewIHostService()
	fileService      = services.NewIFileService()
	groupService     = services.NewIGroupService()
	commandService   = services.NewICommandService()
	settingService   = services.NewISettingService()
	processService   = services.NewIProcessService()
	cronjobService   = services.NewICronjobService()
)
