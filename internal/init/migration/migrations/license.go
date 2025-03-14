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
				LicenseID:  "LIC-TEST-001",
				ExpiryDate: time.Now().AddDate(1, 0, 0),
				IssuedAt:   time.Now().Unix(),
			},
			{
				LicenseID:  "LIC-TEST-002",
				ExpiryDate: time.Now().AddDate(2, 0, 0),
				IssuedAt:   time.Now().Add(-24 * time.Hour).Unix(),
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

var AddLastRemindedAtColumn = &gormigrate.Migration{
	ID: "20240315-add-last-reminded-at",
	Migrate: func(tx *gorm.DB) error {
		// SQLite 需要单独执行ALTER语句
		if tx.Dialector.Name() == "sqlite" {
			return tx.Exec(
				"ALTER TABLE licenses ADD COLUMN last_reminded_at DATETIME",
			).Error
		}
		return tx.AutoMigrate(&models.License{})
	},
	Rollback: func(tx *gorm.DB) error {
		if tx.Dialector.Name() == "sqlite" {
			return tx.Exec(
				"CREATE TEMPORARY TABLE licenses_backup AS SELECT id, created_at, updated_at, deleted_at, license_id, expiry_date, issued_at FROM licenses;" +
					"DROP TABLE licenses;" +
					"CREATE TABLE licenses (" +
					"id INTEGER PRIMARY KEY AUTOINCREMENT," +
					"created_at DATETIME," +
					"updated_at DATETIME," +
					"deleted_at DATETIME," +
					"license_id TEXT," +
					"expiry_date DATETIME," +
					"issued_at INTEGER" +
					");" +
					"INSERT INTO licenses SELECT id, created_at, updated_at, deleted_at, license_id, expiry_date, issued_at FROM licenses_backup;" +
					"DROP TABLE licenses_backup;",
			).Error
		}
		return tx.Migrator().DropColumn(&models.License{}, "last_reminded_at")
	},
}
