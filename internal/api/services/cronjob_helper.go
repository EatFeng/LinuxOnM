package services

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/cmd"
	"fmt"
	"os"
	"strings"
	"time"
)

func (u *CronjobService) HandleJob(cronjob *models.Cronjob) {
	var (
		message []byte
		err     error
	)
	record := cronjobRepo.StartRecords(cronjob.ID, cronjob.KeepLocal, "")
	go func() {
		switch cronjob.Type {
		case "shell":
			if len(cronjob.Script) == 0 {
				return
			}
			record.Records = u.generateLogsPath(*cronjob, record.StartTime)
			_ = cronjobRepo.UpdateRecords(record.ID, map[string]interface{}{"records": record.Records})
			script := cronjob.Script
			if len(cronjob.ContainerName) != 0 {
				command := "sh"
				if len(cronjob.Command) != 0 {
					command = cronjob.Command
				}
				script = fmt.Sprintf("docker exec %s %s -c \"%s\"", cronjob.ContainerName, command, strings.ReplaceAll(cronjob.Script, "\"", "\\\""))
			}
			err = u.handleShell(cronjob.Type, cronjob.Name, script, record.Records)
			u.removeExpiredLog(*cronjob)
		}

		if err != nil {
			if len(message) != 0 {
				record.Records, _ = mkdirAndWriteFile(cronjob, record.StartTime, message)
			}
			cronjobRepo.EndRecords(record, constant.StatusFailed, err.Error(), record.Records)
			return
		}
		if len(message) != 0 {
			record.Records, err = mkdirAndWriteFile(cronjob, record.StartTime, message)
			if err != nil {
				global.LOG.Errorf("save file %s failed, err: %v", record.Records, err)
			}
		}
		cronjobRepo.EndRecords(record, constant.StatusSuccess, "", record.Records)
	}()
}

func (u *CronjobService) generateLogsPath(cronjob models.Cronjob, startTime time.Time) string {
	dir := fmt.Sprintf("%s/task/%s/%s", constant.DataDir, cronjob.Type, cronjob.Name)
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(dir, os.ModePerm)
	}

	path := fmt.Sprintf("%s/%s.log", dir, startTime.Format(constant.DateTimeSlimLayout))
	return path
}

func (u *CronjobService) handleShell(cronType, cornName, script, logPath string) error {
	handleDir := fmt.Sprintf("%s/task/%s/%s", constant.DataDir, cronType, cornName)
	if _, err := os.Stat(handleDir); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(handleDir, os.ModePerm); err != nil {
			return err
		}
	}
	if err := cmd.ExecCronjobWithTimeOut(script, handleDir, logPath, 24*time.Hour); err != nil {
		return err
	}
	return nil
}

func (u *CronjobService) removeExpiredLog(cronjob models.Cronjob) {
	records, _ := cronjobRepo.ListRecord(cronjobRepo.WithByJobID(int(cronjob.ID)), commonRepo.WithOrderBy("created_at desc"))
	if len(records) <= int(cronjob.RetainCopies) {
		return
	}
	for i := int(cronjob.RetainCopies); i < len(records); i++ {
		if len(records[i].File) != 0 {
			files := strings.Split(records[i].File, ",")
			for _, file := range files {
				_ = os.Remove(file)
			}
		}
		_ = cronjobRepo.DeleteRecord(commonRepo.WithByID(uint(records[i].ID)))
		_ = os.Remove(records[i].Records)
	}
}
