package migrations

import (
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableOperationLog = &gormigrate.Migration{
	ID: "20241203-add-table-operation-logs",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&models.OperationLog{}, &models.LoginLog{})
	},
}
