package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type HostRepo struct{}

type IHostRepo interface {
	Get(opts ...DBOption) (models.Host, error)
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
