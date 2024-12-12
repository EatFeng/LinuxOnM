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
	})
	if err := m.Migrate(); err != nil {
		global.LOG.Error(err)
		panic(err)
	}
	global.LOG.Info("Migration run successfully")
}
