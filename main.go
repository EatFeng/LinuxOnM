package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// 创建 Gin 路由
	router := gin.Default()

	// 注册路由
	// router.GET("/", func(c *gin.Context) {
	//     c.JSON(200, gin.H{
	//         "message": "Hello, World!",
	//     })
	// })

	// 启动服务
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
