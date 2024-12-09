package migrations

import (
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableCronjob = &gormigrate.Migration{
	ID: "20241209-add-table-cronjob",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&models.Cronjob{}, &models.JobRecords{})
	},
}
