package request

import "LinuxOnM/internal/utils/files"

type FileOption struct {
	files.FileOption
}

type FileReadByLineReq struct {
	Page     int    `json:"page" validate:"required"`
	PageSize int    `json:"pageSize" validate:"required"`
	Type     string `json:"type" validate:"required"`
	ID       uint   `json:"ID"`
	Name     string `json:"name"`
	Latest   bool   `json:"latest"`
}

type FileCreate struct {
	Path      string `json:"path" validate:"required"`
	Content   string `json:"content"`
	IsDir     bool   `json:"isDir"`
	Mode      int64  `json:"mode"`
	IsLink    bool   `json:"isLink"`
	IsSymlink bool   `json:"isSymlink"`
	LinkPath  string `json:"linkPath"`
	Sub       bool   `json:"sub"`
}
