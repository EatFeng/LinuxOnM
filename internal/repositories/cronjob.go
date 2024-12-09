package repositories

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
	"time"
)

type CronjobRepo struct{}

type ICronjobRepo interface {
	Get(opts ...DBOption) (models.Cronjob, error)
	Create(cronjob *models.Cronjob) error
	WithByJobID(id int) DBOption
	ListRecord(opts ...DBOption) ([]models.JobRecords, error)
	StartRecords(cronjobID uint, fromLocal bool, targetPath string) models.JobRecords
	UpdateRecords(id uint, vars map[string]interface{}) error
	DeleteRecord(opts ...DBOption) error
	EndRecords(record models.JobRecords, status, message, records string)
}

func NewICronjobRepo() ICronjobRepo {
	return &CronjobRepo{}
}

func (u *CronjobRepo) Get(opts ...DBOption) (models.Cronjob, error) {
	var cronjob models.Cronjob
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&cronjob).Error
	return cronjob, err
}

func (u *CronjobRepo) Create(cronjob *models.Cronjob) error {
	return global.DB.Create(cronjob).Error
}

func (c *CronjobRepo) WithByJobID(id int) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("cronjob_id = ?", id)
	}
}

func (u *CronjobRepo) ListRecord(opts ...DBOption) ([]models.JobRecords, error) {
	var cronjobs []models.JobRecords
	db := global.DB.Model(&models.JobRecords{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&cronjobs).Error
	return cronjobs, err
}

func (u *CronjobRepo) StartRecords(cronjobID uint, fromLocal bool, targetPath string) models.JobRecords {
	var record models.JobRecords
	record.StartTime = time.Now()
	record.CronjobID = cronjobID
	record.FromLocal = fromLocal
	record.Status = constant.StatusWaiting
	if err := global.DB.Create(&record).Error; err != nil {
		global.LOG.Errorf("create record status failed, err: %v", err)
	}
	return record
}

func (u *CronjobRepo) UpdateRecords(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.JobRecords{}).Where("id = ?", id).Updates(vars).Error
}

func (u *CronjobRepo) DeleteRecord(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.JobRecords{}).Error
}

func (u *CronjobRepo) EndRecords(record models.JobRecords, status, message, records string) {
	errMap := make(map[string]interface{})
	errMap["records"] = records
	errMap["status"] = status
	errMap["file"] = record.File
	errMap["message"] = message
	errMap["interval"] = time.Since(record.StartTime).Milliseconds()
	if err := global.DB.Model(&models.JobRecords{}).Where("id = ?", record.ID).Updates(errMap).Error; err != nil {
		global.LOG.Errorf("update record status failed, err: %v", err)
	}
}
