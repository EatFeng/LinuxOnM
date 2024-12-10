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
