package storage_client

import (
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/storage_client/client"
)

type StorageClient interface {
	ListBuckets() ([]interface{}, error)
	ListObjects(prefix string) ([]string, error)
	Exist(path string) (bool, error)
	Delete(path string) (bool, error)
	Upload(src, target string) (bool, error)
	Download(src, target string) (bool, error)

	Size(path string) (int64, error)
}

func NewStorageClient(backupType string, vars map[string]interface{}) (StorageClient, error) {
	switch backupType {
	case constant.Local:
		return client.NewLocalClient(vars)
	default:
		return nil, constant.ErrNotSupportType
	}
}
