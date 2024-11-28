package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type LogRepository struct{}

type ILogRepository interface {
	PageLoginLog(limit, offset int, opts ...DBOption) (int64, []models.LoginLog, error)
	CreateLoginLog(user *models.LoginLog) error

	WithByIP(ip string) DBOption
	WithByStatus(status string) DBOption
}

func NewLogRepository() ILogRepository {
	return &LogRepository{}
}

func (u *LogRepository) PageLoginLog(page, size int, opts ...DBOption) (int64, []models.LoginLog, error) {
	var ops []models.LoginLog
	db := global.DB.Model(&models.LoginLog{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&ops).Error
	return count, ops, err
}

func (u *LogRepository) CreateLoginLog(user *models.LoginLog) error {
	return global.DB.Create(user).Error
}

func (u *LogRepository) WithByIP(ip string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(ip) == 0 {
			return db
		}
		return db.Where("ip LIKE ?", "%"+ip+"%")
	}
}

func (u *LogRepository) WithByStatus(status string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		if len(status) == 0 {
			return db
		}
		return db.Where("status = ?", status)
	}
}
