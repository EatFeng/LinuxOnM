package response

import "LinuxOnM/internal/utils/files"

type FileInfo struct {
	files.FileInfo
}

type FileLineContent struct {
	Content string   `json:"content"`
	End     bool     `json:"end"`
	Path    string   `json:"path"`
	Total   int      `json:"total"`
	Lines   []string `json:"lines"`
}

type FileTree struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	IsDir     bool       `json:"isDir"`
	Extension string     `json:"extension"`
	Children  []FileTree `json:"children"`
}
