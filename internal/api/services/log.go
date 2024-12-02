package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/copier"
	"github.com/pkg/errors"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type LogService struct{}

type ILogService interface {
	CreateLoginLog(operation models.LoginLog) error
	PageLoginLog(search dto.SearchLoginLogWithPage) (int64, interface{}, error)

	CreateOperationLog(operation models.OperationLog) error
	PageOperationLog(search dto.SearchOpLogWithPage) (int64, interface{}, error)

	ListSystemLogFile() ([]string, error)

	LoadSSHLog() (string, error)
}

func NewILogService() ILogService {
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

func (u *LogService) ListSystemLogFile() ([]string, error) {
	logDir := path.Join(global.CONF.System.BaseDir, "LinuxOnM/log")
	var files []string
	if err := filepath.Walk(logDir, func(pathItem string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(info.Name(), "LinuxOnM") {
			if info.Name() == "LinuxOnM.log" {
				files = append(files, time.Now().Format("2006-01-02"))
				return nil
			}
			itemFileName := strings.TrimPrefix(info.Name(), "LinuxOnM-")
			itemFileName = strings.TrimSuffix(itemFileName, ".gz")
			itemFileName = strings.TrimSuffix(itemFileName, ".log")
			files = append(files, itemFileName)
			return nil
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if len(files) < 2 {
		return files, nil
	}
	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	return files, nil
}
