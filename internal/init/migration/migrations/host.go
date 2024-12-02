package migrations

import (
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddTableHost = &gormigrate.Migration{
	ID: "20241130-add-table-host",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Host{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.Group{}); err != nil {
			return err
		}
		if err := tx.AutoMigrate(&models.Command{}); err != nil {
			return err
		}
		group := models.Group{
			Name: "default", Type: "host", IsDefault: true,
		}
		if err := tx.Create(&group).Error; err != nil {
			return err
		}
		host := models.Host{
			Name: "localhost", Addr: "127.0.0.1", User: "root", Port: 22, AuthMode: "password", GroupID: group.ID,
		}
		if err := tx.Create(&host).Error; err != nil {
			return err
		}
		return nil
	},
}
