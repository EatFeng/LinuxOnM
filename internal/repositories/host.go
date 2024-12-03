package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type HostRepo struct{}

type IHostRepo interface {
	Get(opts ...DBOption) (models.Host, error)
	WithByAddr(addr string) DBOption
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

func (h *HostRepo) WithByAddr(addr string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("addr = ?", addr)
	}
}
