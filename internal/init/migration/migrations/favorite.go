package migrations

import (
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddFavorite = &gormigrate.Migration{
	ID: "20241213-add-favorite",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Favorite{}); err != nil {
			return err
		}
		return nil
	},
}
