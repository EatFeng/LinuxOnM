package models

import (
	"time"

	"gorm.io/gorm"
)

type License struct {
	gorm.Model
	LicenseID    string    `gorm:"uniqueIndex;size:50"`
	ExpiryDate   time.Time `gorm:"type:date"` // 明确指定为date类型
	HardwareHash string    `gorm:"size:64"`
	IssuedAt     int64     // 确保类型与Python的int一致
}

// 修改过期检查方法
func (l *License) IsExpired() bool {
	return time.Now().UTC().After(l.ExpiryDate.UTC()) // 统一时区处理
}
