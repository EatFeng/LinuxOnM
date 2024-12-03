package services

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/cmd"
	"LinuxOnM/internal/utils/common"
	"LinuxOnM/internal/utils/copier"
	"LinuxOnM/internal/utils/qqwry"
	"fmt"
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

	LoadSSHLog(req dto.SearchSSHLog) (*dto.SSHLog, error)
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

type sshFileItem struct {
	Name string
	Year int
}

func (u *LogService) LoadSSHLog(req dto.SearchSSHLog) (*dto.SSHLog, error) {
	var fileList []sshFileItem
	var data dto.SSHLog
	baseDir := "/var/log"
	if err := filepath.Walk(baseDir, func(pathItem string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (strings.HasPrefix(info.Name(), "secure") || strings.HasPrefix(info.Name(), "auth")) {
			if !strings.HasSuffix(info.Name(), ".gz") {
				fileList = append(fileList, sshFileItem{Name: pathItem, Year: info.ModTime().Year()})
				return nil
			}
			itemFileName := strings.TrimSuffix(pathItem, ".gz")
			if _, err := os.Stat(itemFileName); err != nil && os.IsNotExist(err) {
				if err := handleGunzip(pathItem); err == nil {
					fileList = append(fileList, sshFileItem{Name: itemFileName, Year: info.ModTime().Year()})
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	fileList = sortFileList(fileList)

	command := ""
	if len(req.Info) != 0 {
		command = fmt.Sprintf(" | grep '%s'", req.Info)
	}

	showCountFrom := (req.Page - 1) * req.PageSize
	showCountTo := req.Page * req.PageSize
	nyc, _ := time.LoadLocation(common.LoadTimeZoneByCmd())
	qqWry, err := qqwry.NewQQwry()
	if err != nil {
		global.LOG.Errorf("load qqwry datas failed: %s", err)
	}
	for _, file := range fileList {
		commandItem := ""
		if strings.HasPrefix(path.Base(file.Name), "secure") {
			switch req.Status {
			case constant.StatusSuccess:
				commandItem = fmt.Sprintf("cat %s | grep -a Accepted %s", file.Name, command)
			case constant.StatusFailed:
				commandItem = fmt.Sprintf("cat %s | grep -a 'Failed password for' %s", file.Name, command)
			default:
				commandItem = fmt.Sprintf("cat %s | grep -aE '(Failed password for|Accepted)' %s", file.Name, command)
			}
		}
		if strings.HasPrefix(path.Base(file.Name), "auth.log") {
			switch req.Status {
			case constant.StatusSuccess:
				commandItem = fmt.Sprintf("cat %s | grep -a Accepted %s", file.Name, command)
			case constant.StatusFailed:
				commandItem = fmt.Sprintf("cat %s | grep -aE 'Failed password for|Connection closed by authenticating user' %s", file.Name, command)
			default:
				commandItem = fmt.Sprintf("cat %s | grep -aE \"(Failed password for|Connection closed by authenticating user|Accepted)\" %s", file.Name, command)
			}
		}
		dataItem, successCount, failedCount := loadSSHData(commandItem, showCountFrom, showCountTo, file.Year, qqWry, nyc)
		data.FailedCount += failedCount
		data.TotalCount += successCount + failedCount
		showCountFrom = showCountFrom - (successCount + failedCount)
		showCountTo = showCountTo - (successCount + failedCount)
		data.Logs = append(data.Logs, dataItem...)
	}

	data.SuccessfulCount = data.TotalCount - data.FailedCount
	return &data, nil
}

func sortFileList(fileNames []sshFileItem) []sshFileItem {
	if len(fileNames) < 2 {
		return fileNames
	}
	if strings.HasPrefix(path.Base(fileNames[0].Name), "secure") {
		var itemFile []sshFileItem
		sort.Slice(fileNames, func(i, j int) bool {
			return fileNames[i].Name > fileNames[j].Name
		})
		itemFile = append(itemFile, fileNames[len(fileNames)-1])
		itemFile = append(itemFile, fileNames[:len(fileNames)-1]...)
		return itemFile
	}
	sort.Slice(fileNames, func(i, j int) bool {
		return fileNames[i].Name < fileNames[j].Name
	})
	return fileNames
}

func loadSSHData(command string, showCountFrom, showCountTo, currentYear int, qqWry *qqwry.QQwry, nyc *time.Location) ([]dto.SSHHistory, int, int) {
	var (
		datas        []dto.SSHHistory
		successCount int
		failedCount  int
	)
	stdout2, err := cmd.Exec(command)
	if err != nil {
		return datas, 0, 0
	}
	lines := strings.Split(string(stdout2), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		var itemData dto.SSHHistory
		switch {
		case strings.Contains(lines[i], "Failed password for"):
			itemData = loadFailedSecureDatas(lines[i])
			if len(itemData.Address) != 0 {
				if successCount+failedCount >= showCountFrom && successCount+failedCount < showCountTo {
					itemData.Area = qqWry.Find(itemData.Address).Area
					itemData.Date = loadDate(currentYear, itemData.DateStr, nyc)
					datas = append(datas, itemData)
				}
				failedCount++
			}
		case strings.Contains(lines[i], "Connection closed by authenticating user"):
			itemData = loadFailedAuthDatas(lines[i])
			if len(itemData.Address) != 0 {
				if successCount+failedCount >= showCountFrom && successCount+failedCount < showCountTo {
					itemData.Area = qqWry.Find(itemData.Address).Area
					itemData.Date = loadDate(currentYear, itemData.DateStr, nyc)
					datas = append(datas, itemData)
				}
				failedCount++
			}
		case strings.Contains(lines[i], "Accepted "):
			itemData = loadSuccessDatas(lines[i])
			if len(itemData.Address) != 0 {
				if successCount+failedCount >= showCountFrom && successCount+failedCount < showCountTo {
					itemData.Area = qqWry.Find(itemData.Address).Area
					itemData.Date = loadDate(currentYear, itemData.DateStr, nyc)
					datas = append(datas, itemData)
				}
				successCount++
			}
		}
	}
	return datas, successCount, failedCount
}

func loadFailedSecureDatas(line string) dto.SSHHistory {
	var data dto.SSHHistory
	parts := strings.Fields(line)
	index, dataStr := analyzeDateStr(parts)
	if dataStr == "" {
		return data
	}
	data.DateStr = dataStr
	if strings.Contains(line, " invalid ") {
		data.AuthMode = parts[4+index]
		index += 2
	} else {
		data.AuthMode = parts[4+index]
	}
	data.User = parts[6+index]
	data.Address = parts[8+index]
	data.Port = parts[10+index]
	data.Status = constant.StatusFailed
	if strings.Contains(line, ": ") {
		data.Message = strings.Split(line, ": ")[1]
	}
	return data
}

func loadFailedAuthDatas(line string) dto.SSHHistory {
	var data dto.SSHHistory
	parts := strings.Fields(line)
	index, dataStr := analyzeDateStr(parts)
	if dataStr == "" {
		return data
	}
	data.DateStr = dataStr
	switch index {
	case 1:
		data.User = parts[9]
	case 2:
		data.User = parts[10]
	default:
		data.User = parts[7]
	}
	data.AuthMode = parts[6+index]
	data.Address = parts[9+index]
	data.Port = parts[11+index]
	data.Status = constant.StatusFailed
	if strings.Contains(line, ": ") {
		data.Message = strings.Split(line, ": ")[1]
	}
	return data
}

func loadSuccessDatas(line string) dto.SSHHistory {
	var data dto.SSHHistory
	parts := strings.Fields(line)
	index, dataStr := analyzeDateStr(parts)
	if dataStr == "" {
		return data
	}
	data.DateStr = dataStr
	data.AuthMode = parts[4+index]
	data.User = parts[6+index]
	data.Address = parts[8+index]
	data.Port = parts[10+index]
	data.Status = constant.StatusSuccess
	return data
}

func loadDate(currentYear int, DateStr string, nyc *time.Location) time.Time {
	itemDate, err := time.ParseInLocation("2006 Jan 2 15:04:05", fmt.Sprintf("%d %s", currentYear, DateStr), nyc)
	if err != nil {
		itemDate, _ = time.ParseInLocation("2006 Jan 2 15:04:05", DateStr, nyc)
	}
	return itemDate
}

func analyzeDateStr(parts []string) (int, string) {
	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err == nil {
		if len(parts) < 12 {
			return 0, ""
		}
		return 0, t.Format("2006 Jan 2 15:04:05")
	}
	t, err = time.Parse(constant.DateTimeLayout, fmt.Sprintf("%s %s", parts[0], parts[1]))
	if err == nil {
		if len(parts) < 14 {
			return 0, ""
		}
		return 1, t.Format("2006 Jan 2 15:04:05")
	}

	if len(parts) < 14 {
		return 0, ""
	}
	return 2, fmt.Sprintf("%s %s %s", parts[0], parts[1], parts[2])
}
