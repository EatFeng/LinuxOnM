package db

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"LinuxOnM/internal/global"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init() {
	ensureDbDir(global.CONF.System.DbPath)
	ensureDbFile(global.CONF.System.DbPath, global.CONF.System.DbFile)

	newLogger := createLogger()
	initMonitorDB(newLogger)

	db, err := openDatabase(path.Join(global.CONF.System.DbPath, global.CONF.System.DbFile), newLogger)
	if err != nil {
		panic(err)
	}
	configureDatabase(db)

	global.DB = db
	global.LOG.Info("init db successfully")
}

func initMonitorDB(newLogger logger.Interface) {
	ensureDbDir(global.CONF.System.DbPath)
	ensureDbFile(global.CONF.System.DbPath, "monitor.db")

	db, err := openDatabase(path.Join(global.CONF.System.DbPath, "monitor.db"), newLogger)
	if err != nil {
		panic(err)
	}
	configureDatabase(db)

	global.MonitorDB = db
	global.LOG.Info("init monitor db successfully")
}

func ensureDbDir(dir string) {
	if _, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(fmt.Errorf("init db dir failed, err: %v", err))
		}
	}
}

func ensureDbFile(dir, file string) {
	fullPath := path.Join(dir, file)
	if _, err := os.Stat(fullPath); err != nil {
		f, err := os.Create(fullPath)
		if err != nil {
			panic(fmt.Errorf("init db file failed, err: %v", err))
		}
		_ = f.Close()
	}
}

func createLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

func openDatabase(fullPath string, newLogger logger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(fullPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func configureDatabase(db *gorm.DB) {
	_ = db.Exec("PRAGMA journal_mode = WAL;")
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
