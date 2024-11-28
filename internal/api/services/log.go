package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/copier"
	"github.com/pkg/errors"
)

type LogService struct{}

type ILogService interface {
	CreateLoginLog(operation models.LoginLog) error
	PageLoginLog(search dto.SearchLoginLogWithPage) (int64, interface{}, error)

	CreateOperationLog(operation models.OperationLog)
	PageOperationLog(search dto.SearchOpLogWithPage) (int64, interface{}, error)
}

func NewLogService() ILogService {
	return &LogService{}
}

func (u *LogService) PageLoginLog(req dto.SearchLoginLogWithPage) (int64, interface{}, error) {
	total, ops, err := logRepo.PageLoginLog(
		req.Page,
		req.PageSize,
		logRepo.WithByIP(req.IP),
		logRepo.WithByStatus(req.Status),
		commonRepo.WithOrderBy("created_at DESC"),
	)
	var dtoOps []dto.LoginLog
	for _, op := range ops {
		var item dto.LoginLog
		if err := copier.Copy(&item, op); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		dtoOps = append(dtoOps, item)
	}
	return total, dtoOps, err
}

func (u *LogService) CreateLoginLog(operation models.LoginLog) error {
	return logRepo.CreateLoginLog(&operation)
}

func (u *LogService) PageOperationLog(req dto.SearchOpLogWithPage) (int64, interface{}, error) {
	total, ops, err := logRepo.PageOperationLog(
		req.Page,
		req.PageSize,
		logRepo.WithByGroup(req.Source),
		logRepo.WithByLikeOperation(req.Operation),
		logRepo.WithByStatus(req.Status),
		commonRepo.WithOrderBy("created_at desc"),
	)
	// transform models.OperationLog to dto.OperationLog
	var dtoOps []dto.OperationLog
	for _, op := range ops {
		var item dto.OperationLog
		if err := copier.Copy(&item, &op); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		dtoOps = append(dtoOps, item)
	}
	return total, dtoOps, err
}

func (u *LogService) CreateOperationLog(operation models.OperationLog) error {
	return logRepo.CreateOperationLog(&operation)
}
