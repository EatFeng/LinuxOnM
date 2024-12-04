package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/copier"
	"github.com/pkg/errors"
)

type GroupService struct{}

type IGroupService interface {
	List(req dto.GroupSearch) ([]dto.GroupInfo, error)
	Create(req dto.GroupCreate) error
}

func NewIGroupService() IGroupService {
	return &GroupService{}
}

func (u *GroupService) List(req dto.GroupSearch) ([]dto.GroupInfo, error) {
	groups, err := groupRepo.GetList(commonRepo.WithByType(req.Type), commonRepo.WithOrderBy("is_default desc"), commonRepo.WithOrderBy("created_at desc"))
	if err != nil {
		return nil, constant.ErrRecordNotFound
	}
	var dtoUsers []dto.GroupInfo
	for _, group := range groups {
		var item dto.GroupInfo
		if err := copier.Copy(&item, &group); err != nil {
			return nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		dtoUsers = append(dtoUsers, item)
	}
	return dtoUsers, err
}

func (u *GroupService) Create(req dto.GroupCreate) error {
	group, _ := groupRepo.Get(commonRepo.WithByName(req.Name), commonRepo.WithByType(req.Type))
	if group.ID != 0 {
		return constant.ErrRecordExist
	}
	if err := copier.Copy(&group, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	if err := groupRepo.Create(&group); err != nil {
		return err
	}
	return nil
}
