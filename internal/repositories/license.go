package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"time"
)

type LicenseRepo struct{}

type ILicenseRepo interface {
	Create(license *models.License) error
	// GetLatestValid(hardwareHash string) (*models.License, error)
	GetLatestValid() (*models.License, error)
}

func NewLicenseRepo() ILicenseRepo {
	return &LicenseRepo{}
}

// 创建许可证
func (r *LicenseRepo) Create(license *models.License) error {
	return global.DB.Create(license).Error
}

// 获取最新的有效许可证
// func (r *LicenseRepo) GetLatestValid(hardwareHash string) (*models.License, error) {
func (r *LicenseRepo) GetLatestValid() (*models.License, error) {
	var license models.License
	err := global.DB.Where("expiry_date > ?", time.Now().UTC()).Order("issued_at DESC").First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}
