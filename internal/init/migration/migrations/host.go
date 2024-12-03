package migrations

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/encrypt"
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
			Name: "localhost", Addr: "127.0.0.1", User: "feng-yite", Port: 22, AuthMode: "password",
			GroupID: group.ID, Password: "qwer1234",
		}
		if err := tx.Create(&host).Error; err != nil {
			return err
		}
		return nil
	},
}

var EncryptHostPassword = &gormigrate.Migration{
	ID: "20241203-encrypt-host-password",
	Migrate: func(tx *gorm.DB) error {
		var hosts []models.Host
		if err := tx.Where("1 = 1").Find(&hosts).Error; err != nil {
			return err
		}

		var encryptSetting models.Setting
		if err := tx.Where("key = ?", "EncryptKey").Find(&encryptSetting).Error; err != nil {
			return err
		}
		global.CONF.System.EncryptKey = encryptSetting.Value

		for _, host := range hosts {
			if len(host.Password) != 0 {
				pass, err := encrypt.StringEncrypt(host.Password)
				if err != nil {
					return err
				}
				if err := tx.Model(&models.Host{}).Where("id = ?", host.ID).Update("password", pass).Error; err != nil {
					return err
				}
			}
			if len(host.PrivateKey) != 0 {
				key, err := encrypt.StringEncrypt(host.PrivateKey)
				if err != nil {
					return err
				}
				if err := tx.Model(&models.Host{}).Where("id = ?", host.ID).Update("private_key", key).Error; err != nil {
					return err
				}
			}
			if len(host.PassPhrase) != 0 {
				pass, err := encrypt.StringEncrypt(host.PassPhrase)
				if err != nil {
					return err
				}
				if err := tx.Model(&models.Host{}).Where("id = ?", host.ID).Update("pass_phrase", pass).Error; err != nil {
					return err
				}
			}
		}
		return nil
	},
}
