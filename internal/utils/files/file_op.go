package files

import (
	"github.com/spf13/afero"
)

type FileOp struct {
	Fs afero.Fs
}

type Process struct {
	Total   uint64  `json:"total"`
	Written uint64  `json:"written"`
	Percent float64 `json:"percent"`
	Name    string  `json:"name"`
}

func NewFileOp() FileOp {
	return FileOp{
		Fs: afero.NewOsFs(),
	}
}

func (f FileOp) Stat(dst string) bool {
	info, _ := f.Fs.Stat(dst)
	return info != nil
}
