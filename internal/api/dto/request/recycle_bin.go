package request

type RecycleBinCreate struct {
	SourcePath string `json:"sourcePath" validate:"required"`
}
