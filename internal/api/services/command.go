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
	Delete(ids []uint) error
	Update(id uint, upMap map[string]interface{}) error
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

func (u *CommandService) Delete(ids []uint) error {
	if len(ids) == 1 {
		command, _ := commandRepo.Get(commonRepo.WithByID(ids[0]))
		if command.ID == 0 {
			return constant.ErrRecordNotFound
		}
		return commandRepo.Delete(commonRepo.WithByID(ids[0]))
	}
	return commandRepo.Delete(commonRepo.WithIDsIn(ids))
}

func (u *CommandService) Update(id uint, upMap map[string]interface{}) error {
	return commandRepo.Update(id, upMap)
}
