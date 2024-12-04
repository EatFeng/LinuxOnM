package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type GroupRepo struct{}

type IGroupRepo interface {
	Get(opts ...DBOption) (models.Group, error)
	GetList(opts ...DBOption) ([]models.Group, error)
	WithByHostDefault() DBOption
}

func NewIGroupRepo() IGroupRepo {
	return &GroupRepo{}
}

func (u *GroupRepo) Get(opts ...DBOption) (models.Group, error) {
	var group models.Group
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&group).Error
	return group, err
}

func (u *GroupRepo) GetList(opts ...DBOption) ([]models.Group, error) {
	var groups []models.Group
	db := global.DB.Model(&models.Group{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&groups).Error
	return groups, err
}

func (u *GroupRepo) WithByHostDefault() DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("is_default = ? AND type = ?", 1, "host")
	}
}
