package routers

import (
	"LinuxOnM/internal/api/handlers"
	"LinuxOnM/internal/middleware"
	"github.com/gin-gonic/gin"
)

type TerminalRouter struct{}

func (s *TerminalRouter) InitRouter(Router *gin.RouterGroup) {
	terminalRouter := Router.Group("terminal").Use(middleware.PasswordExpired())

	baseApi := handlers.ApiGroupApp.BaseApi
	{
		terminalRouter.GET("", baseApi.WsSsh)
	}
}
