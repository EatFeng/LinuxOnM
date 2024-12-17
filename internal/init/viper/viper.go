package viper

import (
	"LinuxOnM/cmd/conf"
	"LinuxOnM/internal/configs"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/files"

	"fmt"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
	"path"

	"github.com/spf13/viper"
)

func Init() {
	baseDir := "/opt"
	bindaddress := "0.0.0.0"
	port := "9999"
	mode := ""
	version := "v1.0.0"
	username, password, entrance := "", "", ""
	dbfile := "main.db"
	fileOp := files.NewFileOp()
	v := viper.NewWithOptions()
	v.SetConfigType("yaml")

	config := configs.ServerConfig{}
	if err := yaml.Unmarshal(conf.AppYaml, &config); err != nil {
		panic(err)
	}
	if config.System.Mode != "" {
		mode = config.System.Mode
	}
	if mode == "dev" && fileOp.Stat("/opt/LinuxOnM/conf/app.yaml") {
		v.SetConfigName("app")
		v.AddConfigPath(path.Join("/opt/LinuxOnM/conf"))
		if err := v.ReadInConfig(); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
	v.OnConfigChange(func(e fsnotify.Event) {
		if err := v.Unmarshal(&global.CONF); err != nil {
			panic(err)
		}
	})
	serverConfig := configs.ServerConfig{}
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	if mode == "dev" && fileOp.Stat("/opt/LinuxOnM/conf/app.yaml") {
		if serverConfig.System.BaseDir != "" {
			baseDir = serverConfig.System.BaseDir
		}
		if serverConfig.System.BindAddress != "" {
			bindaddress = serverConfig.System.BindAddress
		}
		if serverConfig.System.Port != "" {
			port = serverConfig.System.Port
		}
		if serverConfig.System.Version != "" {
			version = serverConfig.System.Version
		}
		if serverConfig.System.Username != "" {
			username = serverConfig.System.Username
		}
		if serverConfig.System.Password != "" {
			password = serverConfig.System.Password
		}
		if serverConfig.System.Entrance != "" {
			entrance = serverConfig.System.Entrance
		}
		if serverConfig.System.DbFile != "" {
			dbfile = serverConfig.System.DbFile
		}
	}

	global.CONF = serverConfig
	global.CONF.System.BaseDir = baseDir
	global.CONF.System.BindAddress = bindaddress
	global.CONF.System.Port = port
	global.CONF.System.Version = version
	global.CONF.System.Username = username
	global.CONF.System.Password = password
	global.CONF.System.Entrance = entrance
	global.Viper = v
	global.CONF.System.DataDir = path.Join(global.CONF.System.BaseDir, "LinuxOnM")
	global.CONF.System.DbPath = path.Join(global.CONF.System.DataDir, "db")
	global.CONF.System.DbFile = dbfile
	global.CONF.System.LogPath = path.Join(global.CONF.System.DataDir, "log")
	global.CONF.System.Cache = path.Join(global.CONF.System.DataDir, "cache")
	global.CONF.System.Backup = path.Join(global.CONF.System.DataDir, "backup")
	global.CONF.System.TmpDir = path.Join(global.CONF.System.DataDir, "tmp")
}
