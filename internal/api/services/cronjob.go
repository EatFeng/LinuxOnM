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
	Delete(req dto.CronjobBatchDelete) error
	Update(id uint, req dto.CronjobUpdate) error
	HandleOnce(id uint) error
	UpdateStatus(id uint, status string) error
	SearchWithPage(search dto.PageCronjob) (int64, interface{}, error)

	StartJob(cronjob *models.Cronjob, isUpdate bool) (string, error)
	SearchRecords(search dto.SearchRecord) (int64, interface{}, error)
	LoadRecordLog(req dto.OperateByID) string
	CleanRecord(req dto.CronjobClean) error
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

func (u *CronjobService) Delete(req dto.CronjobBatchDelete) error {
	for _, id := range req.IDs {
		cronjob, _ := cronjobRepo.Get(commonRepo.WithByID(id))
		if cronjob.ID == 0 {
			return errors.New("find cronjob in db failed")
		}
		ids := strings.Split(cronjob.EntryIDs, ",")
		for _, id := range ids {
			idItem, _ := strconv.Atoi(id)
			global.Cron.Remove(cron.EntryID(idItem))
		}
		global.LOG.Infof("stop cronjob entryID: %s", cronjob.EntryIDs)
		if err := u.CleanRecord(dto.CronjobClean{CronjobID: id, CleanData: req.CleanData, IsDelete: true}); err != nil {
			return err
		}
		if err := cronjobRepo.Delete(commonRepo.WithByID(id)); err != nil {
			return err
		}
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

func (u *CronjobService) Update(id uint, req dto.CronjobUpdate) error {
	var cronjob models.Cronjob
	if err := copier.Copy(&cronjob, &req); err != nil {
		return errors.WithMessage(constant.ErrStructTransform, err.Error())
	}
	cronModel, err := cronjobRepo.Get(commonRepo.WithByID(id))
	if err != nil {
		return constant.ErrRecordNotFound
	}
	upMap := make(map[string]interface{})
	cronjob.EntryIDs = cronModel.EntryIDs
	cronjob.Type = cronModel.Type
	spec := cronjob.Spec
	if cronModel.Status == constant.StatusEnable {
		newEntryIDs, err := u.StartJob(&cronjob, true)
		if err != nil {
			return err
		}
		upMap["entry_ids"] = newEntryIDs
	} else {
		ids := strings.Split(cronjob.EntryIDs, ",")
		for _, id := range ids {
			idItem, _ := strconv.Atoi(id)
			global.Cron.Remove(cron.EntryID(idItem))
		}
	}

	upMap["name"] = req.Name
	upMap["spec"] = spec
	upMap["script"] = req.Script
	upMap["command"] = req.Command
	upMap["container_name"] = req.ContainerName
	upMap["exclusion_rules"] = req.ExclusionRules
	upMap["url"] = req.URL
	upMap["source_dir"] = req.SourceDir

	upMap["backup_accounts"] = req.BackupAccounts
	upMap["default_download"] = req.DefaultDownload
	upMap["retain_copies"] = req.RetainCopies
	upMap["secret"] = req.Secret
	return cronjobRepo.Update(id, upMap)
}

func (u *CronjobService) UpdateStatus(id uint, status string) error {
	cronjob, _ := cronjobRepo.Get(commonRepo.WithByID(id))
	if cronjob.ID == 0 {
		return errors.WithMessage(constant.ErrRecordNotFound, "record not found")
	}
	var (
		entryIDs string
		err      error
	)
	if status == constant.StatusEnable {
		entryIDs, err = u.StartJob(&cronjob, false)
		if err != nil {
			return err
		}
	} else {
		ids := strings.Split(cronjob.EntryIDs, ",")
		for _, id := range ids {
			idItem, _ := strconv.Atoi(id)
			global.Cron.Remove(cron.EntryID(idItem))
		}
		global.LOG.Infof("stop cronjob entryID: %s", cronjob.EntryIDs)
	}
	return cronjobRepo.Update(cronjob.ID, map[string]interface{}{"status": status, "entry_ids": entryIDs})
}

func (u *CronjobService) HandleOnce(id uint) error {
	cronjob, _ := cronjobRepo.Get(commonRepo.WithByID(id))
	if cronjob.ID == 0 {
		return constant.ErrRecordNotFound
	}
	u.HandleJob(&cronjob)
	return nil
}

func (u *CronjobService) SearchRecords(search dto.SearchRecord) (int64, interface{}, error) {
	total, records, err := cronjobRepo.PageRecords(
		search.Page,
		search.PageSize,
		commonRepo.WithByStatus(search.Status),
		cronjobRepo.WithByJobID(search.CronjobID),
		commonRepo.WithByDate(search.StartTime, search.EndTime))
	var dtoCronjobs []dto.Record
	for _, record := range records {
		var item dto.Record
		if err := copier.Copy(&item, &record); err != nil {
			return 0, nil, errors.WithMessage(constant.ErrStructTransform, err.Error())
		}
		item.StartTime = record.StartTime.Format(constant.DateTimeLayout)
		dtoCronjobs = append(dtoCronjobs, item)
	}
	return total, dtoCronjobs, err
}

func (u *CronjobService) LoadRecordLog(req dto.OperateByID) string {
	record, err := cronjobRepo.GetRecord(commonRepo.WithByID(req.ID))
	if err != nil {
		return ""
	}
	if _, err := os.Stat(record.Records); err != nil {
		return ""
	}
	content, err := os.ReadFile(record.Records)
	if err != nil {
		return ""
	}
	return string(content)
}

func (u *CronjobService) CleanRecord(req dto.CronjobClean) error {
	cronjob, err := cronjobRepo.Get(commonRepo.WithByID(req.CronjobID))
	if err != nil {
		return err
	}
	if req.CleanData {
		if hasBackup(cronjob.Type) {
			accountMap, err := loadClientMap(cronjob.BackupAccounts)
			if err != nil {
				return err
			}
			cronjob.RetainCopies = 0
			u.removeExpiredBackup(cronjob, accountMap, models.BackupRecord{})
		} else {
			u.removeExpiredLog(cronjob)
		}
	}
	if req.IsDelete {
		records, _ := backupRepo.ListRecord(backupRepo.WithByCronID(cronjob.ID))
		for _, records := range records {
			records.CronjobID = 0
			_ = backupRepo.UpdateRecord(&records)
		}
	}
	delRecords, err := cronjobRepo.ListRecord(cronjobRepo.WithByJobID(int(req.CronjobID)))
	if err != nil {
		return err
	}
	for _, del := range delRecords {
		_ = os.RemoveAll(del.Records)
	}
	if err := cronjobRepo.DeleteRecord(cronjobRepo.WithByJobID(int(req.CronjobID))); err != nil {
		return err
	}
	return nil
}
