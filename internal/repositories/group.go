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
	Create(group *models.Group) error
	Update(id uint, vars map[string]interface{}) error
	Delete(opts ...DBOption) error
	CancelDefault(groupType string) error
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

func (u *GroupRepo) Create(group *models.Group) error {
	return global.DB.Create(group).Error
}

func (u *GroupRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.Group{}).Where("id = ?", id).Updates(vars).Error
}

func (u *GroupRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.Group{}).Error
}

func (u *GroupRepo) WithByHostDefault() DBOption {
	return func(g *gorm.DB) *gorm.DB {
		return g.Where("is_default = ? AND type = ?", 1, "host")
	}
}

func (u *GroupRepo) CancelDefault(groupType string) error {
	return global.DB.Model(&models.Group{}).Where("is_default = ? AND type = ?", 1, groupType).Updates(map[string]interface{}{"is_default": 0}).Error
}
