package global

import (
	"LinuxOnM/internal/configs"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	CONF      configs.ServerConfig
	Viper     *viper.Viper
	LOG       *logrus.Logger
	DB        *gorm.DB
	MonitorDB *gorm.DB
	VILID     *validator.Validate
)
