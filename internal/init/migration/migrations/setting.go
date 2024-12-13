package migrations

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/common"
	"LinuxOnM/internal/utils/encrypt"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
	"time"
)

var AddTableSetting = &gormigrate.Migration{
	ID: "20241126_add_table_setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&models.Setting{}); err != nil {
			return err
		}
		encryptKey := common.RandStr(16)
		if err := tx.Create(&models.Setting{Key: "UserName", Value: global.CONF.System.Username}).Error; err != nil {
			return err
		}
		global.CONF.System.EncryptKey = encryptKey
		pass, _ := encrypt.StringEncrypt(global.CONF.System.Password)
		if err := tx.Create(&models.Setting{Key: "Password", Value: pass}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "SecurityEntrance", Value: global.CONF.System.Entrance}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "SessionTimeout", Value: "86400"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "LocalTime", Value: ""}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "ServerPort", Value: global.CONF.System.Port}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "JWTSigningKey", Value: common.RandStr(16)}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "EncryptKey", Value: encryptKey}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "ExpirationTime", Value: time.Now().AddDate(0, 0, 10).Format(constant.DateTimeLayout)}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "ExpirationDays", Value: "0"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "ComplexityVerification", Value: "enable"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "MonitorStatus", Value: "enable"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "MonitorStoreDays", Value: "7"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "MessageType", Value: "none"}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "SystemVersion", Value: global.CONF.System.Version}).Error; err != nil {
			return err
		}

		if err := tx.Create(&models.Setting{Key: "SystemStatus", Value: "Free"}).Error; err != nil {
			return err
		}

		return nil
	},
}

var AddBindAndAllowIPs = &gormigrate.Migration{
	ID: "20241205-add-cronjob-settings",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "TimeZone", Value: common.LoadTimeZoneByCmd()}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "NtpSite", Value: "pool.ntp.org"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "MonitorInterval", Value: "1"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddProxy = &gormigrate.Migration{
	ID: "20241212-add-proxy-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "ProxyType", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "ProxyUrl", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "ProxyPort", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "ProxyUser", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "ProxyPasswd", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "ProxyPasswdKeep", Value: ""}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddBindAddress = &gormigrate.Migration{
	ID: "20241212-add-bind-address",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "BindAddress", Value: "0.0.0.0"}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "Ipv6", Value: "disable"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddPanelName = &gormigrate.Migration{
	ID: "20241213-add-panelname-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "PanelName", Value: "LinuxOnM"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddDefaultNetwork = &gormigrate.Migration{
	ID: "20241213-add-default-network-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "DefaultNetwork", Value: "ens33"}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddLaskInfo = &gormigrate.Migration{
	ID: "20241213-add-last-info-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "LastCleanTime", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "LastCleanSize", Value: ""}).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Setting{Key: "LastCleanData", Value: ""}).Error; err != nil {
			return err
		}
		return nil
	},
}

var AddTheme = &gormigrate.Migration{
	ID: "20241213-add-theme-setting",
	Migrate: func(tx *gorm.DB) error {
		if err := tx.Create(&models.Setting{Key: "Theme", Value: "light"}).Error; err != nil {
			return err
		}
		return nil
	},
}
