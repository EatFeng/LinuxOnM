package services

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type NotificationService struct {
	APIURL string
}

func NewNotificationService() *NotificationService {
	// 从配置获取通知URL（示例从系统配置获取）
	urlSetting, _ := settingRepo.Get(settingRepo.WithByKey("NotificationURL"))
	return &NotificationService{
		APIURL: urlSetting.Value,
	}
}

// SendAlert 发送报警通知（线程安全）
func (s *NotificationService) SendAlert(metricType string, currentValue float64, durationSeconds int) {
	if s.APIURL == "" {
		global.LOG.Error("The notification API is not configured. Skipping alarm sending.")
		return
	}

	// build alarm data
	data := s.buildNotificationData(metricType, currentValue, durationSeconds)

	// 异步发送防止阻塞
	go func() {
		jsonData, _ := json.Marshal(data)
		resp, err := http.Post(s.APIURL, "application/json", bytes.NewBuffer(jsonData))

		if err != nil {
			global.LOG.Errorf("通知发送失败: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			global.LOG.Errorf("通知接口返回错误状态码: %d", resp.StatusCode)
		}
	}()
}

// 格式化报警描述（英文）
func (s *NotificationService) formatDescription(metric string, value float64, duration int) string {
	return fmt.Sprintf("%s %.2f%% exceeds threshold for %d seconds",
		metric, value, duration)
}

var (
	idGenerator = models.NewAlarmIDGenerator()
)

func (s *NotificationService) buildNotificationData(metric string, value float64, duration int) models.NotificationData {
	var eventCode, description string

	switch metric {
	case "CPU":
		eventCode = models.EventCodeCPUHighUsage
		description = "CPU usage"
	case "Memory":
		eventCode = models.EventCodeMemoryHigeUsage
		description = "Memory usage"
	default:
		eventCode = models.EventCodeUnknown
		description = "Warning Unknown"
	}

	// 生成唯一报警码（示例：ALARM-1710779323-CPU-0001)
	alarmID := idGenerator.Next(metric)

	// build alarm description
	alarmDetail := s.formatDescription(description, value, duration)

	return models.NotificationData{
		EventCode: eventCode,
		AlarmTime: time.Now().Format("2006-01-02 15:04:05"),
		DevNumber: alarmID,
		DevType:   alarmDetail,
	}
}
