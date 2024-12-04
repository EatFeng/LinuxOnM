package repositories

import (
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
)

type CommandRepo struct{}

type ICommandRepo interface {
	GetList(opts ...DBOption) ([]models.Command, error)
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
