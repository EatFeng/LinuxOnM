package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/copier"
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"os"
	"strconv"
	"strings"
	"time"
)

type CronjobService struct{}

type ICronjobService interface {
	Create(cronjobDto dto.CronjobCreate) error
	SearchWithPage(search dto.PageCronjob) (int64, interface{}, error)

	StartJob(cronjob *models.Cronjob, isUpdate bool) (string, error)
}

func NewICronjobService() ICronjobService {
	return &CronjobService{}
}

func (u *CronjobService) Create(cronjobDto dto.CronjobCreate) error {
	cronjob, _ := cronjobRepo.Get(commonRepo.WithByName(cronjobDto.Name))
	if cronjob.ID != 0 {
		return constant.ErrRecordExist
	}
	cronjob.Secret = cronjobDto.Secret
	if err := copier.Copy(&cronjob, &cronjobDto); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	cronjob.Status = constant.StatusEnable

	global.LOG.Infof("create cronjob %s successful, spec: %s", cronjob.Name, cronjob.Spec)
	spec := cronjob.Spec
	entryIDs, err := u.StartJob(&cronjob, false)
	if err != nil {
		return err
	}
	cronjob.Spec = spec
	cronjob.EntryIDs = entryIDs
	if err := cronjobRepo.Create(&cronjob); err != nil {
		return err
	}
	return nil
}

func (u *CronjobService) StartJob(cronjob *models.Cronjob, isUpdate bool) (string, error) {
	if len(cronjob.EntryIDs) != 0 && isUpdate {
		ids := strings.Split(cronjob.EntryIDs, ",")
		for _, id := range ids {
			idItem, _ := strconv.Atoi(id)
			global.Cron.Remove(cron.EntryID(idItem))
		}
	}
	specs := strings.Split(cronjob.Spec, ",")
	var ids []string
	for _, spec := range specs {
		cronjob.Spec = spec
		entryID, err := u.AddCronJob(cronjob)
		if err != nil {
			return "", err
		}
		ids = append(ids, fmt.Sprintf("%v", entryID))
	}
	return strings.Join(ids, ","), nil
}

func (u *CronjobService) AddCronJob(cronjob *models.Cronjob) (int, error) {
	addFunc := func() {
		u.HandleJob(cronjob)
	}
	global.LOG.Infof("add %s job %s successful", cronjob.Type, cronjob.Name)
	entryID, err := global.Cron.AddFunc(cronjob.Spec, addFunc)
	if err != nil {
		return 0, err
	}
	global.LOG.Infof("start cronjob entryID: %d", entryID)
	return int(entryID), nil
}

func mkdirAndWriteFile(cronjob *models.Cronjob, startTime time.Time, msg []byte) (string, error) {
	dir := fmt.Sprintf("%s/task/%s/%s", constant.DataDir, cronjob.Type, cronjob.Name)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}
	}

	path := fmt.Sprintf("%s/%s.log", dir, startTime.Format(constant.DateTimeSlimLayout))
	global.LOG.Infof("cronjob %s has generated some logs %s", cronjob.Name, path)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	_, _ = write.WriteString(string(msg))
	write.Flush()
	return path, nil
}

func (u *CronjobService) SearchWithPage(search dto.PageCronjob) (int64, interface{}, error) {
	total, cronjobs, err := cronjobRepo.Page(search.Page, search.PageSize, commonRepo.WithLikeName(search.Info), commonRepo.WithOrderRuleBy(search.OrderBy, search.Order))
	var dtoCronjobs []dto.CronjobInfo
	for _, cronjob := range cronjobs {
		var item dto.CronjobInfo
		if err := copier.Copy(&item, &cronjob); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		record, _ := cronjobRepo.RecordFirst(cronjob.ID)
		if record.ID != 0 {
			item.LastRecordTime = record.StartTime.Format(constant.DateTimeLayout)
		} else {
			item.LastRecordTime = "-"
		}
		dtoCronjobs = append(dtoCronjobs, item)
	}
	return total, dtoCronjobs, err
}
