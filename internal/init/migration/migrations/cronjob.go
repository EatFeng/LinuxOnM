package migrations

import (
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddNewTableCronjob = &gormigrate.Migration{
	ID: "20241209-add-new-table-cronjob-two",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&models.Cronjob{}, &models.JobRecords{})
	},
}
