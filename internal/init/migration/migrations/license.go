package migrations

import (
	"LinuxOnM/internal/models"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableLicense = &gormigrate.Migration{
	ID: "20250312-add-table-license",
	Migrate: func(tx *gorm.DB) error {
		// 创建许可证表
		if err := tx.AutoMigrate(&models.License{}); err != nil {
			return err
		}

		// 添加测试数据（可选）
		testLicenses := []models.License{
			{
				LicenseID:    "LIC-TEST-001",
				ExpiryDate:   time.Now().AddDate(1, 0, 0),
				HardwareHash: "test_hash_001",
				IssuedAt:     time.Now().Unix(),
			},
			{
				LicenseID:    "LIC-TEST-002",
				ExpiryDate:   time.Now().AddDate(2, 0, 0),
				HardwareHash: "test_hash_002",
				IssuedAt:     time.Now().Add(-24 * time.Hour).Unix(),
			},
		}
		return tx.Create(&testLicenses).Error
	},
	Rollback: func(tx *gorm.DB) error {
		// 删除索引
		if err := tx.Exec("DROP INDEX IF EXISTS idx_license_id").Error; err != nil {
			return err
		}
		// 删除数据表
		return tx.Migrator().DropTable(&models.License{})
	},
}
