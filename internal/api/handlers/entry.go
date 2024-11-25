package handlers

import "LinuxOnM/internal/api/services"

type BaseApi struct{}

type ApiGroup struct {
	BaseApi
}

var ApiGroupApp = new(ApiGroup)

var (
	dashboardService = services.NewDashboardService()
)
