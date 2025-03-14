package dto

import "time"

type CronjobCreate struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
	Spec string `json:"spec" validate:"required"`

	Script         string `json:"script"`
	Command        string `json:"command"`
	ContainerName  string `json:"containerName"`
	ExclusionRules string `json:"exclusionRules"`
	SourceDir      string `json:"sourceDir"`

	BackupAccounts  string `json:"backupAccounts"`
	DefaultDownload string `json:"defaultDownload"`
	RetainCopies    int    `json:"retainCopies" validate:"number,min=1"`
	Secret          string `json:"secret"`
}

type PageCronjob struct {
	PageInfo
	Info    string `json:"info"`
	OrderBy string `json:"orderBy" validate:"required,oneof=name status created_at"`
	Order   string `json:"order" validate:"required,oneof=null ascending descending"`
}

type CronjobInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Spec string `json:"spec"`

	Script          string `json:"script"`
	Command         string `json:"command"`
	ContainerName   string `json:"containerName"`
	ExclusionRules  string `json:"exclusionRules"`
	URL             string `json:"url"`
	SourceDir       string `json:"sourceDir"`
	BackupAccounts  string `json:"backupAccounts"`
	DefaultDownload string `json:"defaultDownload"`
	RetainCopies    int    `json:"retainCopies"`

	LastRecordTime string `json:"lastRecordTime"`
	Status         string `json:"status"`
	Secret         string `json:"secret"`
}

type CronjobUpdate struct {
	ID   uint   `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
	Spec string `json:"spec" validate:"required"`

	Script         string `json:"script"`
	Command        string `json:"command"`
	ContainerName  string `json:"containerName"`
	ExclusionRules string `json:"exclusionRules"`
	URL            string `json:"url"`
	SourceDir      string `json:"sourceDir"`

	BackupAccounts  string `json:"backupAccounts"`
	DefaultDownload string `json:"defaultDownload"`
	RetainCopies    int    `json:"retainCopies" validate:"number,min=1"`
	Secret          string `json:"secret"`
}

type CronjobUpdateStatus struct {
	ID     uint   `json:"id" validate:"required"`
	Status string `json:"status" validate:"required"`
}

type CronjobBatchDelete struct {
	CleanData bool   `json:"cleanData"`
	IDs       []uint `json:"ids" validate:"required"`
}

type CronjobClean struct {
	IsDelete  bool `json:"isDelete"`
	CleanData bool `json:"cleanData"`
	CronjobID uint `json:"cronjobID" validate:"required"`
}

type SearchRecord struct {
	PageInfo
	CronjobID int       `json:"cronjobID"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
	Status    string    `json:"status"`
}

type Record struct {
	ID         uint   `json:"id"`
	StartTime  string `json:"startTime"`
	Records    string `json:"records"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	TargetPath string `json:"targetPath"`
	Interval   int    `json:"interval"`
	File       string `json:"file"`
}
