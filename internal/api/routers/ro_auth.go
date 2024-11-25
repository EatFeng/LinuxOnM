package routers

import (
	"LinuxOnM/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

type AuthRouter struct{}

func (s *AuthRouter) InitRouter(Router *gin.RouterGroup) {
	authRouter := Router.Group("auth")
	baseApi := handlers.ApiGroupApp.BaseApi
	{
		authRouter.GET("/login", baseApi.Login)
	}
}
