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
	SearchForTree() ([]dto.CommandTree, error)
	Create(commandDto dto.CommandOperate) error
	Delete(ids []uint) error
	Update(id uint, upMap map[string]interface{}) error
	SearchWithPage(search dto.SearchCommandWithPage) (int64, interface{}, error)
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

func (u *CommandService) SearchForTree() ([]dto.CommandTree, error) {
	cmdList, err := commandRepo.GetList(commonRepo.WithOrderBy("name"))
	if err != nil {
		return nil, err
	}
	groups, err := groupRepo.GetList(commonRepo.WithByType("command"))
	if err != nil {
		return nil, err
	}
	var lists []dto.CommandTree
	for _, group := range groups {
		var data dto.CommandTree
		data.ID = group.ID + 10000
		data.Label = group.Name
		for _, cmd := range cmdList {
			if cmd.GroupID == group.ID {
				data.Children = append(data.Children, dto.CommandInfo{ID: cmd.ID, Name: cmd.Name, Command: cmd.Command})
			}
		}
		if len(data.Children) != 0 {
			lists = append(lists, data)
		}
	}
	return lists, err
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

func (u *CommandService) SearchWithPage(search dto.SearchCommandWithPage) (int64, interface{}, error) {
	total, commands, err := commandRepo.Page(search.Page, search.PageSize, commandRepo.WithLikeName(search.Name), commonRepo.WithLikeName(search.Info), commonRepo.WithByGroupID(search.GroupID), commonRepo.WithOrderRuleBy(search.OrderBy, search.Order))
	if err != nil {
		return 0, nil, err
	}
	groups, _ := groupRepo.GetList(commonRepo.WithByType("command"), commonRepo.WithOrderBy("name"))
	var dtoCommands []dto.CommandInfo
	for _, command := range commands {
		var item dto.CommandInfo
		if err := copier.Copy(&item, &command); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		for _, group := range groups {
			if command.GroupID == group.ID {
				item.GroupBelong = group.Name
				item.GroupID = group.ID
			}
		}
		dtoCommands = append(dtoCommands, item)
	}
	return total, dtoCommands, err
}
