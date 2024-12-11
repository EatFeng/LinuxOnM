package services

import (
	"LinuxOnM/internal/models"
	"LinuxOnM/internal/utils/storage_client"
	"encoding/json"
)

type BackupService struct{}

type IBackupService interface {
	NewClient(backup *models.BackupAccount) (storage_client.StorageClient, error)
}

func NewIBackupService() IBackupService {
	return &BackupService{}
}

func (u *BackupService) NewClient(backup *models.BackupAccount) (storage_client.StorageClient, error) {
	varMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(backup.Vars), &varMap); err != nil {
		return nil, err
	}
	varMap["bucket"] = backup.Bucket

	backClient, err := storage_client.NewStorageClient(backup.Type, varMap)
	if err != nil {
		return nil, err
	}

	return backClient, nil
}
