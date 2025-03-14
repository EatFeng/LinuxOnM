package models

import (
	"time"

	"gorm.io/gorm"
)

type License struct {
	gorm.Model
	LicenseID      string    `gorm:"uniqueIndex;size:50"`
	ExpiryDate     time.Time `gorm:"type:datetime"`
	IssuedAt       int64
	LastRemindedAt time.Time `gorm:"type:datetime"`
}

func (l *License) IsExpired() bool {
	now := time.Now().UTC()
	expiry := l.ExpiryDate.UTC()

	// 精确到秒级判断
	return now.After(expiry) || now.Equal(expiry)
}
