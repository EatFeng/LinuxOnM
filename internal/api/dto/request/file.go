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

type FileDelete struct {
	Path        string `json:"path" validate:"required"`
	IsDir       bool   `json:"isDir"`
	ForceDelete bool   `json:"forceDelete"`
}

type FileContentReq struct {
	Path     string `json:"path" validate:"required"`
	IsDetail bool   `json:"isDetail"`
}

type FileRoleReq struct {
	Paths []string `json:"paths" validate:"required"`
	Mode  int64    `json:"mode" validate:"required"`
	User  string   `json:"user" validate:"required"`
	Group string   `json:"group" validate:"required"`
	Sub   bool     `json:"sub"`
}

type FileCompress struct {
	Files   []string `json:"files" validate:"required"`
	Dst     string   `json:"dst" validate:"required"`
	Type    string   `json:"type" validate:"required"`
	Name    string   `json:"name" validate:"required"`
	Replace bool     `json:"replace"`
	Secret  string   `json:"secret"`
}

type FilePathCheck struct {
	Path string `json:"path" validate:"required"`
}

type FileMove struct {
	Type     string   `json:"type" validate:"required"`
	OldPaths []string `json:"oldPaths" validate:"required"`
	NewPath  string   `json:"newPath" validate:"required"`
	Name     string   `json:"name"`
	Cover    bool     `json:"cover"`
}

type FileRename struct {
	OldName string `json:"oldName" validate:"required"`
	NewName string `json:"newName" validate:"required"`
}

type FileDeCompress struct {
	Dst    string `json:"dst"  validate:"required"`
	Type   string `json:"type"  validate:"required"`
	Path   string `json:"path" validate:"required"`
	Secret string `json:"secret"`
}
