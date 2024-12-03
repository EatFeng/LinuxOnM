package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type HostRepo struct{}

type IHostRepo interface {
	Get(opts ...DBOption) (models.Host, error)
	Create(host *models.Host) error
	Update(id uint, vars map[string]interface{}) error

	WithByAddr(addr string) DBOption
	WithByUser(user string) DBOption
	WithByPort(port uint) DBOption
}

func NewIHostRepo() IHostRepo {
	return &HostRepo{}
}

func (h *HostRepo) Get(opts ...DBOption) (models.Host, error) {
	var host models.Host
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&host).Error
	return host, err
}

func (h *HostRepo) Create(host *models.Host) error {
	return global.DB.Create(host).Error
}

func (h *HostRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.Host{}).Where("id = ?", id).Updates(vars).Error
}

func (h *HostRepo) WithByAddr(addr string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("addr = ?", addr)
	}
}

func (h *HostRepo) WithByUser(user string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("user = ?", user)
	}
}

func (h *HostRepo) WithByPort(port uint) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("port = ?", port)
	}
}
