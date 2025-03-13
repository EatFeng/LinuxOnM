package middleware

import (
	"LinuxOnM/internal/repositories"
	"LinuxOnM/internal/utils/encrypt"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LicenseCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("license check")

		// 1. 生成硬件哈希
		hardwareHash, err := encrypt.GenerateHardwareHash()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "无法生成硬件指纹",
			})
			return
		}

		// 2. 查询最新有效许可证
		licenseRepo := repositories.NewLicenseRepo()
		license, err := licenseRepo.GetLatestValid(hardwareHash)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "未找到有效许可证",
			})
			return
		}

		// 3. 检查是否过期
		if license.IsExpired() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "许可证已过期",
			})
			return
		}

		// 4. 注入许可证信息到上下文
		c.Set("license", license)
		c.Next()
	}
}
