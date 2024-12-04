package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/copier"
	"github.com/pkg/errors"
)

type CommandService struct{}

type ICommandService interface {
	List() ([]dto.CommandInfo, error)
	Create(commandDto dto.CommandOperate) error
}

func NewICommandService() ICommandService {
	return &CommandService{}
}

func (u *CommandService) List() ([]dto.CommandInfo, error) {
	commands, err := commandRepo.GetList(commonRepo.WithOrderBy("name"))
	if err != nil {
		return nil, constant.ErrRecordNotFound
	}
	var dtoCommands []dto.CommandInfo
	for _, command := range commands {
		var item dto.CommandInfo
		if err := copier.Copy(&item, &command); err != nil {
			return nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		dtoCommands = append(dtoCommands, item)
	}
	return dtoCommands, err
}

func (u *CommandService) Create(commandDto dto.CommandOperate) error {
	command, _ := commandRepo.Get(commonRepo.WithByName(commandDto.Name))
	if command.ID != 0 {
		return constant.ErrRecordExist
	}
	if err := copier.Copy(&command, &commandDto); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	if err := commandRepo.Create(&command); err != nil {
		return err
	}
	return nil
}
