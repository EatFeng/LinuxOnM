package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type SettingRepo struct{}

type ISettingRepo interface {
	Get(opts ...DBOption) (models.Setting, error)

	WithByKey(key string) DBOption
}

func NewISettingRepo() ISettingRepo {
	return &SettingRepo{}
}

func (u *SettingRepo) Get(opts ...DBOption) (models.Setting, error) {
	var settings models.Setting
	db := global.DB.Model(&models.Setting{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&settings).Error
	return settings, err
}

func (c *SettingRepo) WithByKey(key string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("key = ?", key)
	}
}
