package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"gorm.io/gorm"
)

type CommandRepo struct{}

type ICommandRepo interface {
	GetList(opts ...DBOption) ([]models.Command, error)
	Get(opts ...DBOption) (models.Command, error)
	Create(command *models.Command) error
	Delete(opts ...DBOption) error
	Update(id uint, vars map[string]interface{}) error
	Page(limit, offset int, opts ...DBOption) (int64, []models.Command, error)
	WithLikeName(name string) DBOption
}

func NewICommandRepo() ICommandRepo {
	return &CommandRepo{}
}

func (u *CommandRepo) GetList(opts ...DBOption) ([]models.Command, error) {
	var commands []models.Command
	db := global.DB.Model(&models.Command{})
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.Find(&commands).Error
	return commands, err
}

func (u *CommandRepo) Get(opts ...DBOption) (models.Command, error) {
	var command models.Command
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	err := db.First(&command).Error
	return command, err
}

func (u *CommandRepo) Create(command *models.Command) error {
	return global.DB.Create(command).Error
}

func (u *CommandRepo) Delete(opts ...DBOption) error {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db.Delete(&models.Command{}).Error
}

func (u *CommandRepo) Update(id uint, vars map[string]interface{}) error {
	return global.DB.Model(&models.Command{}).Where("id = ?", id).Updates(vars).Error
}

func (u *CommandRepo) Page(page, size int, opts ...DBOption) (int64, []models.Command, error) {
	var users []models.Command
	db := global.DB.Model(&models.Command{})
	for _, opt := range opts {
		db = opt(db)
	}
	count := int64(0)
	db = db.Count(&count)
	err := db.Limit(size).Offset(size * (page - 1)).Find(&users).Error
	return count, users, err
}

func (a CommandRepo) WithLikeName(name string) DBOption {
	return func(g *gorm.DB) *gorm.DB {
		if len(name) == 0 {
			return g
		}
		return g.Where("name like ? or command like ?", "%"+name+"%", "%"+name+"%")
	}
}
