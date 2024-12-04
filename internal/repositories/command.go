package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type CommandRepo struct{}

type ICommandRepo interface {
	GetList(opts ...DBOption) ([]models.Command, error)
	Get(opts ...DBOption) (models.Command, error)
	Create(command *models.Command) error
	Delete(opts ...DBOption) error
	Update(id uint, vars map[string]interface{}) error
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
