package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
	"time"
)

type SettingRepo struct{}

type ISettingRepo interface {
	Get(opts ...DBOption) (models.Setting, error)
	GetList(opts ...DBOption) ([]models.Setting, error)
	Update(key, value string) error

	WithByKey(key string) DBOption

	CreateMonitorBase(model models.MonitorBase) error
	BatchCreateMonitorIO(ioList []models.MonitorIO) error
	BatchCreateMonitorNet(ioList []models.MonitorNetwork) error
	DelMonitorBase(timeForDelete time.Time) error
	DelMonitorIO(timeForDelete time.Time) error
	DelMonitorNet(timeForDelete time.Time) error
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

func (u *SettingRepo) GetList(opts ...DBOption) ([]models.Setting, error) {
	var settings []models.Setting
	db := global.DB.Model(&models.Setting{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&settings).Error
	return settings, err
}

func (u *SettingRepo) Update(key, value string) error {
	return global.DB.Model(&models.Setting{}).Where("key = ?", key).Updates(map[string]interface{}{"value": value}).Error
}

func (c *SettingRepo) WithByKey(key string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("key = ?", key)
	}
}

func (u *SettingRepo) CreateMonitorBase(model models.MonitorBase) error {
	return global.MonitorDB.Create(&model).Error
}

func (u *SettingRepo) BatchCreateMonitorIO(ioList []models.MonitorIO) error {
	return global.MonitorDB.CreateInBatches(ioList, len(ioList)).Error
}

func (u *SettingRepo) BatchCreateMonitorNet(ioList []models.MonitorNetwork) error {
	return global.MonitorDB.CreateInBatches(ioList, len(ioList)).Error
}

func (u *SettingRepo) DelMonitorBase(timeForDelete time.Time) error {
	return global.MonitorDB.Where("created_at < ?", timeForDelete).Delete(&models.MonitorBase{}).Error
}

func (u *SettingRepo) DelMonitorIO(timeForDelete time.Time) error {
	return global.MonitorDB.Where("created_at < ?", timeForDelete).Delete(&models.MonitorIO{}).Error
}

func (u *SettingRepo) DelMonitorNet(timeForDelete time.Time) error {
	return global.MonitorDB.Where("created_at < ?", timeForDelete).Delete(&models.MonitorNetwork{}).Error
}
