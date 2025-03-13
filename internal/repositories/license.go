package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type LicenseRepo struct{}

type ILicenseRepo interface {
	Create(license *models.License) error
}

func NewLicenseRepo() ILicenseRepo {
	return &LicenseRepo{}
}

// 创建许可证
func (r *LicenseRepo) Create(license *models.License) error {
	return global.DB.Create(license).Error
}
