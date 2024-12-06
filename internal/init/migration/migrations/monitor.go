package migrations

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var NewMonitorDB = &gormigrate.Migration{
	ID: "20240408-new-monitor-db",
	Migrate: func(tx *gorm.DB) error {
		var (
			bases    []models.MonitorBase
			ios      []models.MonitorIO
			networks []models.MonitorNetwork
		)
		_ = tx.Find(&bases).Error
		_ = tx.Find(&ios).Error
		_ = tx.Find(&networks).Error

		if err := global.MonitorDB.AutoMigrate(&models.MonitorBase{}, &models.MonitorIO{}, &models.MonitorNetwork{}); err != nil {
			return err
		}
		_ = global.MonitorDB.Exec("DELETE FROM monitor_bases").Error
		_ = global.MonitorDB.Exec("DELETE FROM monitor_ios").Error
		_ = global.MonitorDB.Exec("DELETE FROM monitor_networks").Error

		if len(bases) != 0 {
			for i := 0; i <= len(bases)/200; i++ {
				var itemData []models.MonitorBase
				if 200*(i+1) <= len(bases) {
					itemData = bases[200*i : 200*(i+1)]
				} else {
					itemData = bases[200*i:]
				}
				if len(itemData) != 0 {
					if err := global.MonitorDB.Create(&itemData).Error; err != nil {
						return err
					}
				}
			}
		}
		if len(ios) != 0 {
			for i := 0; i <= len(ios)/200; i++ {
				var itemData []models.MonitorIO
				if 200*(i+1) <= len(ios) {
					itemData = ios[200*i : 200*(i+1)]
				} else {
					itemData = ios[200*i:]
				}
				if len(itemData) != 0 {
					if err := global.MonitorDB.Create(&itemData).Error; err != nil {
						return err
					}
				}
			}
		}
		if len(networks) != 0 {
			for i := 0; i <= len(networks)/200; i++ {
				var itemData []models.MonitorNetwork
				if 200*(i+1) <= len(networks) {
					itemData = networks[200*i : 200*(i+1)]
				} else {
					itemData = networks[200*i:]
				}
				if len(itemData) != 0 {
					if err := global.MonitorDB.Create(&itemData).Error; err != nil {
						return err
					}
				}
			}
		}
		return nil
	},
}
