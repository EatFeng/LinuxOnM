package dto

type CronjobCreate struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
	Spec string `json:"spec" validate:"required"`

	Script         string `json:"script"`
	Command        string `json:"command"`
	ContainerName  string `json:"containerName"`
	ExclusionRules string `json:"exclusionRules"`
	URL            string `json:"url"`
	SourceDir      string `json:"sourceDir"`

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
	DefaultDownload string `json:"defaultDownload"`
	RetainCopies    int    `json:"retainCopies"`

	LastRecordTime string `json:"lastRecordTime"`
	Status         string `json:"status"`
	Secret         string `json:"secret"`
}
