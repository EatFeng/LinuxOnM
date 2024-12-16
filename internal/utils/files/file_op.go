package files

import (
	"LinuxOnM/internal/utils/cmd"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"io/fs"
	"os"
	"path/filepath"
	"time"
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

func (f FileOp) CreateDir(dst string, mode fs.FileMode) error {
	return f.Fs.MkdirAll(dst, mode)
}

func (f FileOp) CreateDirWithMode(dst string, mode fs.FileMode) error {
	if err := f.Fs.MkdirAll(dst, mode); err != nil {
		return err
	}
	return f.ChmodRWithMode(dst, mode, true)
}

func (f FileOp) ChmodRWithMode(dst string, mode fs.FileMode, sub bool) error {
	cmdStr := fmt.Sprintf(`chmod %v "%s"`, fmt.Sprintf("%o", mode.Perm()), dst)
	if sub {
		cmdStr = fmt.Sprintf(`chmod -R %v "%s"`, fmt.Sprintf("%o", mode.Perm()), dst)
	}
	if cmd.HasNoPasswordSudo() {
		cmdStr = fmt.Sprintf("sudo %s", cmdStr)
	}
	if msg, err := cmd.ExecWithTimeOut(cmdStr, 10*time.Second); err != nil {
		if msg != "" {
			return errors.New(msg)
		}
		return err
	}
	return nil
}

func (f FileOp) LinkFile(source string, dst string, isSymlink bool) error {
	if isSymlink {
		osFs := afero.OsFs{}
		return osFs.SymlinkIfPossible(source, dst)
	} else {
		return os.Link(source, dst)
	}
}

func (f FileOp) CreateFileWithMode(dst string, mode fs.FileMode) error {
	file, err := f.Fs.OpenFile(dst, os.O_CREATE, mode)
	if err != nil {
		return err
	}
	return file.Close()
}

func (f FileOp) DeleteDir(dst string) error {
	return f.Fs.RemoveAll(dst)
}

func (f FileOp) DeleteFile(dst string) error {
	return f.Fs.Remove(dst)
}

func (f FileOp) OpenFile(dst string) (fs.File, error) {
	return f.Fs.Open(dst)
}

func (f FileOp) GetDirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return float64(size), nil
}

func (f FileOp) Mv(oldPath, dstPath string) error {
	cmdStr := fmt.Sprintf(`mv '%s' '%s'`, oldPath, dstPath)
	if err := cmd.ExecCmd(cmdStr); err != nil {
		return err
	}
	return nil
}

func (f FileOp) ChownR(dst string, uid string, gid string, sub bool) error {
	cmdStr := fmt.Sprintf(`chown %s:%s "%s"`, uid, gid, dst)
	if sub {
		cmdStr = fmt.Sprintf(`chown -R %s:%s "%s"`, uid, gid, dst)
	}
	if cmd.HasNoPasswordSudo() {
		cmdStr = fmt.Sprintf("sudo %s", cmdStr)
	}
	if msg, err := cmd.ExecWithTimeOut(cmdStr, 10*time.Second); err != nil {
		if msg != "" {
			return errors.New(msg)
		}
		return err
	}
	return nil
}

func (f FileOp) ChmodR(dst string, mode int64, sub bool) error {
	cmdStr := fmt.Sprintf(`chmod %v "%s"`, fmt.Sprintf("%04o", mode), dst)
	if sub {
		cmdStr = fmt.Sprintf(`chmod -R %v "%s"`, fmt.Sprintf("%04o", mode), dst)
	}
	if cmd.HasNoPasswordSudo() {
		cmdStr = fmt.Sprintf("sudo %s", cmdStr)
	}
	if msg, err := cmd.ExecWithTimeOut(cmdStr, 10*time.Second); err != nil {
		if msg != "" {
			return errors.New(msg)
		}
		return err
	}
	return nil
}
