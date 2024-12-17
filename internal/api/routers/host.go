package routers

import (
	handler "LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"
	"github.com/gin-gonic/gin"
)

type HostRouter struct{}

func (s *HostRouter) InitRouter(Router *gin.RouterGroup) {
	hostRouter := Router.Group("host").Use(middleware.PasswordExpired())
	baseApi := handler.ApiGroupApp.BaseApi
	{
		// host-terminal-terminal & host
		hostRouter.POST("", baseApi.CreateHost)
		hostRouter.POST("/del", baseApi.DeleteHost)
		hostRouter.POST("/search", baseApi.SearchHost)
		hostRouter.POST("/update", baseApi.UpdateHost)
		hostRouter.POST("/update/group", baseApi.UpdateHostGroup)
		hostRouter.POST("/tree", baseApi.HostTree)
		hostRouter.POST("/test/byid/:id", baseApi.TestByID)
		hostRouter.POST("/test/byinfo", baseApi.TestByInfo)
		// host-terminal-command
		hostRouter.GET("/command", baseApi.ListCommand)
		hostRouter.POST("/command", baseApi.CreateCommand)
		hostRouter.POST("/command/del", baseApi.DeleteCommand)
		hostRouter.POST("/command/update", baseApi.UpdateCommand)
		hostRouter.POST("/command/search", baseApi.SearchCommand)
		hostRouter.GET("/command/tree", baseApi.SearchCommandTree)
		// host-monitor
		hostRouter.POST("/monitor/search", baseApi.LoadMonitor)
		hostRouter.POST("/monitor/clean", baseApi.CleanMonitor)
		hostRouter.GET("/monitor/net_options", baseApi.GetNetworkOptions)
		hostRouter.GET("/monitor/io_options", baseApi.GetIOOptions)
		// host-firewall
		hostRouter.GET("/firewall/base", baseApi.LoadFirewallBaseInfo)
		hostRouter.POST("/firewall/operate", baseApi.OperateFirewall)
		hostRouter.POST("/firewall/search", baseApi.SearchFirewallRule)
		hostRouter.POST("/firewall/port", baseApi.OperatePortRule)
		hostRouter.POST("/firewall/forward", baseApi.OperateForwardRule)
		hostRouter.POST("/firewall/ip", baseApi.OperateIPRule)
		hostRouter.POST("/firewall/update/port", baseApi.UpdatePortRule)
		hostRouter.POST("/firewall/update/addr", baseApi.UpdateAddrRule)
		hostRouter.POST("/firewall/update/description", baseApi.UpdateFirewallDescription)
		hostRouter.POST("/firewall/batch", baseApi.BatchOperateRule)
	}
}
