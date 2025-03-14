package middleware

import (
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/repositories"
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	alertEndpoint = "http://上级平台/api/v1/license-expiry-alerts"
	warningDays   = 7
	timeLayout    = "2006-01-02T15:04:05Z"
)

func LicenseCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("license check")

		// 1. 查询最新有效许可证
		licenseRepo := repositories.NewLicenseRepo()
		license, err := licenseRepo.GetLatestValid()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "未找到有效许可证",
			})
			return
		}

		// 2. 检查是否过期
		if license.IsExpired() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "许可证已过期",
			})
			return
		}

		// 3. 临近过期提醒逻辑
		now := time.Now().UTC()
		expiresAt := license.ExpiryDate.UTC()

		// 计算基于UTC零点的准确天数
		remainingDays := calculateRemainingDays(now, expiresAt)

		if remainingDays > 0 && remainingDays <= warningDays {
			handleLicenseWarning(c, license, remainingDays, now)
		}

		// 4. 注入许可证信息到上下文
		c.Set("license", license)
		c.Next()
	}
}

func calculateRemainingDays(now, expiresAt time.Time) int {
	// 对齐到UTC零点
	nowMidnight := now.Truncate(24 * time.Hour)
	expiresMidnight := expiresAt.Truncate(24 * time.Hour)

	// 计算精确天数
	duration := expiresMidnight.Sub(nowMidnight)
	return int(math.Floor(duration.Hours() / 24))
}

func handleLicenseWarning(c *gin.Context, license *models.License, days int, checkTime time.Time) {
	licenseRepo := repositories.NewLicenseRepo()

	// 原子更新提醒时间
	if err := licenseRepo.UpdateLastRemindedAt(license.LicenseID, checkTime); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("[License] 提醒时间更新失败: %v", err)
		}
		return
	}

	// 添加响应头
	c.Header("License-Warning", fmt.Sprintf("License will be invalid after %d day(s)", days))

	// 异步发送报警（防止阻塞请求）
	// go func() {
	// 	alertID := uuid.New().String()
	// 	payload := map[string]interface{}{
	// 		"alert_id":       alertID,
	// 		"license_id":     license.LicenseID,
	// 		"remaining_days": days,
	// 		"expiry_date":    license.ExpiryDate.Format(timeLayout),
	// 		"alert_time":     time.Now().UTC().Format(timeLayout),
	// 	}

	// 	if err := sendLicenseAlert(payload); err != nil {
	// 		log.Printf("[License] 报警推送失败: %v", err)
	// 	}
	// }()
}

// func sendLicenseAlert(payload map[string]interface{}) error {
// 	jsonData, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Errorf("序列化失败: %w", err)
// 	}

// 	req, err := http.NewRequest(http.MethodPost, alertEndpoint, bytes.NewReader(jsonData))
// 	if err != nil {
// 		return fmt.Errorf("创建请求失败: %w", err)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{Timeout: 5 * time.Second}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Errorf("请求发送失败: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode >= http.StatusBadRequest {
// 		return fmt.Errorf("异常响应码: %d", resp.StatusCode)
// 	}

// 	return nil
// }
