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
