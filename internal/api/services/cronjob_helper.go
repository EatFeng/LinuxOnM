package services

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/repositories"
	"LinuxOnM/internal/utils/cmd"
	"LinuxOnM/internal/utils/ntp"
	"LinuxOnM/internal/utils/storage_client"
	"context"
	"fmt"
	"os"
	"path"
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
		case "ntp":
			err = u.handleNtpSync()
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

func hasBackup(cronjobType string) bool {
	return cronjobType == "directory" || cronjobType == "snapshot" || cronjobType == "log"
}

func loadClientMap(backupAccounts string) (map[string]cronjobUploadHelper, error) {
	clients := make(map[string]cronjobUploadHelper)
	accounts, err := backupRepo.List()
	if err != nil {
		return nil, err
	}
	targets := strings.Split(backupAccounts, ",")
	for _, target := range targets {
		if len(target) == 0 {
			continue
		}
		for _, account := range accounts {
			if target == account.Type {
				client, err := NewIBackupService().NewClient(&account)
				if err != nil {
					return nil, err
				}
				pathItem := account.BackupPath
				if account.BackupPath != "/" {
					pathItem = strings.TrimPrefix(account.BackupPath, "/")
				}
				clients[target] = cronjobUploadHelper{
					client:     client,
					backupPath: pathItem,
					backType:   account.Type,
				}
			}
		}
	}
	return clients, nil
}

func (u *CronjobService) removeExpiredBackup(cronjob models.Cronjob, accountMap map[string]cronjobUploadHelper, record models.BackupRecord) {
	var opts []repositories.DBOption
	opts = append(opts, commonRepo.WithByFrom("cronjob"))
	opts = append(opts, backupRepo.WithByCronID(cronjob.ID))
	opts = append(opts, commonRepo.WithOrderBy("created_at desc"))
	if record.ID != 0 {
		opts = append(opts, backupRepo.WithByType(record.Type))
		opts = append(opts, commonRepo.WithByName(record.Name))
		opts = append(opts, backupRepo.WithByDetailName(record.DetailName))
	}
	records, _ := backupRepo.ListRecord(opts...)
	if len(records) <= int(cronjob.RetainCopies) {
		return
	}
	for i := int(cronjob.RetainCopies); i < len(records); i++ {
		accounts := strings.Split(cronjob.BackupAccounts, ",")
		if cronjob.Type == "snapshot" {
			for _, account := range accounts {
				if len(account) != 0 {
					_, _ = accountMap[account].client.Delete(path.Join(accountMap[account].backupPath, "system_snapshot", records[i].FileName))
				}
			}
			_ = snapshotRepo.Delete(commonRepo.WithByName(strings.TrimSuffix(records[i].FileName, ".tar.gz")))
		} else {
			for _, account := range accounts {
				if len(account) != 0 {
					_, _ = accountMap[account].client.Delete(path.Join(accountMap[account].backupPath, records[i].FileDir, records[i].FileName))
				}
			}
		}
		_ = backupRepo.DeleteRecord(context.Background(), commonRepo.WithByID(records[i].ID))
	}
}

func (u *CronjobService) handleNtpSync() error {
	ntpServer, err := settingRepo.Get(settingRepo.WithByKey("NtpSite"))
	if err != nil {
		return err
	}
	ntime, err := ntp.GetRemoteTime(ntpServer.Value)
	if err != nil {
		return err
	}
	if err := ntp.UpdateSystemTime(ntime.Format(constant.DateTimeLayout)); err != nil {
		return err
	}
	return nil
}

type cronjobUploadHelper struct {
	backupPath string
	backType   string
	client     storage_client.StorageClient
}
