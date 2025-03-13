package migration

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/init/migration/migrations"

	"github.com/go-gormigrate/gormigrate/v2"
)

func Init() {
	m := gormigrate.New(global.DB, gormigrate.DefaultOptions, []*gormigrate.Migration{
		migrations.AddTableSetting,
		migrations.AddTableOperationLog,
		migrations.AddTableHost,
		migrations.EncryptHostPassword,
		migrations.NewMonitorDB,
		migrations.AddBindAndAllowIPs,
		migrations.AddNewTableCronjob,
		migrations.AddProxy,
		migrations.AddBindAddress,
		migrations.AddPanelName,
		migrations.AddDefaultNetwork,
		migrations.AddLaskInfo,
		migrations.AddTheme,
		migrations.AddNoAuthSetting,
		migrations.AddFavorite,
		migrations.AddWebsiteCA,
		migrations.AddDefaultCA,
		migrations.AddWebsiteSSL,
		migrations.AddSSLSetting,
		migrations.AddTableFirewall,
		migrations.AddTableForward,
		migrations.AddTableImageRepo,
		migrations.AddTableLicense,
	})
	if err := m.Migrate(); err != nil {
		global.LOG.Error(err)
		panic(err)
	}
	global.LOG.Info("Migration run successfully")
}
