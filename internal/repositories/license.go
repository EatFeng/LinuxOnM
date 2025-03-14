package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"database/sql"
	"time"
)

type LicenseRepo struct{}

type ILicenseRepo interface {
	Create(license *models.License) error
	GetLatestValid() (*models.License, error)
	UpdateLastRemindedAt(licenseID string, remindedAt time.Time) error
}

func NewLicenseRepo() ILicenseRepo {
	return &LicenseRepo{}
}

// 创建许可证
func (r *LicenseRepo) Create(license *models.License) error {
	return global.DB.Create(license).Error
}

// 获取最新的有效许可证
func (r *LicenseRepo) GetLatestValid() (*models.License, error) {
	var license models.License
	err := global.DB.Order("issued_at DESC").First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

func (r *LicenseRepo) UpdateLastRemindedAt(licenseID string, checkTime time.Time) error {
	dayStart := checkTime.UTC().Truncate(24 * time.Hour)

	result := global.DB.Model(&models.License{}).
		Where("license_id = ? AND (last_reminded_at IS NULL OR last_reminded_at < ?)",
			licenseID, dayStart).
		Updates(map[string]interface{}{
			"last_reminded_at": dayStart,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
