package models

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

type NotificationData struct {
	EventCode string `json:"envet_code"`
	AlarmTime string `json:"alarm_time"`
	DevNumber string `json:"dev_number"` // 报警唯一标识码
	DevType   string `json:"dev_type"`   // 报警详情描述（新定义）
}

// 报警唯一码生成器
type AlarmIDGenerator struct {
	counter uint64
	prefix  string
}

func NewAlarmIDGenerator() *AlarmIDGenerator {
	return &AlarmIDGenerator{
		prefix: fmt.Sprintf("ALARM-%d-", time.Now().Unix()),
	}
}

func (g *AlarmIDGenerator) Next(metricType string) string {
	atomic.AddUint64(&g.counter, 1)
	return fmt.Sprintf("%s%s-%04d",
		g.prefix,
		strings.ToUpper(metricType),
		atomic.LoadUint64(&g.counter))
}

// constant event code
const (
	EventCodeCPUHighUsage    = "OP000"
	EventCodeMemoryHigeUsage = "OP111"
	EventCodeUnknown         = "OP999"
)
