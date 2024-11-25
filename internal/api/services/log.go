package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/copier"
	"github.com/pkg/errors"
)

type LogService struct{}

type ILogService interface {
	PageLoginLog(search dto.SearchLoginLogWithPage) (int64, interface{}, error)
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
