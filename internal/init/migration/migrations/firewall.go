package migrations

import (
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableFirewall = &gormigrate.Migration{
	ID: "20241217-add-table-firewall",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Firewall{}); err != nil {
			return err
		}
		return nil
	},
}

var AddTableForward = &gormigrate.Migration{
	ID: "20241217-add-forward",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Forward{}); err != nil {
			return err
		}
		return nil
	},
}
