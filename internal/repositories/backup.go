package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"context"
	"gorm.io/gorm"
)

type BackupRepo struct{}

type IBackupRepo interface {
	List(opts ...DBOption) ([]models.BackupAccount, error)
	ListRecord(opts ...DBOption) ([]models.BackupRecord, error)
	UpdateRecord(record *models.BackupRecord) error
	WithByCronID(cronjobID uint) DBOption
	WithByType(backupType string) DBOption
	WithByDetailName(detailName string) DBOption
	DeleteRecord(ctx context.Context, opts ...DBOption) error
}

func NewIBackupRepo() IBackupRepo {
	return &BackupRepo{}
}

func (u *BackupRepo) List(opts ...DBOption) ([]models.BackupAccount, error) {
	var ops []models.BackupAccount
	db := global.DB.Model(&models.BackupAccount{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&ops).Error
	return ops, err
}

func (u *BackupRepo) WithByCronID(cronjobID uint) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("cronjob_id = ?", cronjobID)
	}
}

func (u *BackupRepo) ListRecord(opts ...DBOption) ([]models.BackupRecord, error) {
	var users []models.BackupRecord
	db := global.DB.Model(&models.BackupRecord{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&users).Error
	return users, err
}

func (u *BackupRepo) UpdateRecord(record *models.BackupRecord) error {
	return global.DB.Save(record).Error
}

func (u *BackupRepo) WithByType(backupType string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(backupType) == 0 {
			return g
		}
		return g.Where("type = ?", backupType)
	}
}

func (u *BackupRepo) WithByDetailName(detailName string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(detailName) == 0 {
			return g
		}
		return g.Where("detail_name = ?", detailName)
	}
}

func (u *BackupRepo) DeleteRecord(ctx context.Context, opts ...DBOption) error {
	return getTx(ctx, opts...).Delete(&models.BackupRecord{}).Error
}
