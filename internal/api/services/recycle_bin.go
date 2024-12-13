package services

import (
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/files"
	"fmt"
	"github.com/shirou/gopsutil/v3/disk"
	"path"
	"strings"
	"time"
)

type RecycleBinService struct {
}

type IRecycleBinService interface {
	Create(create request.RecycleBinCreate) error
}

func NewIRecycleBinService() IRecycleBinService {
	return &RecycleBinService{}
}

func (r RecycleBinService) Create(create request.RecycleBinCreate) error {
	op := files.NewFileOp()
	if !op.Stat(create.SourcePath) {
		return buserr.New(constant.ErrLinkPathNotFound)
	}
	clashDir, err := getClashDir(create.SourcePath)
	if err != nil {
		return err
	}
	paths := strings.Split(create.SourcePath, "/")
	rNamePre := strings.Join(paths, "_LinuxOnM_")
	deleteTime := time.Now()
	openFile, err := op.OpenFile(create.SourcePath)
	if err != nil {
		return err
	}
	fileInfo, err := openFile.Stat()
	if err != nil {
		return err
	}
	size := 0
	if fileInfo.IsDir() {
		sizeF, err := op.GetDirSize(create.SourcePath)
		if err != nil {
			return err
		}
		size = int(sizeF)
	} else {
		size = int(fileInfo.Size())
	}

	rName := fmt.Sprintf("_LinuxOnM_%s%s_p_%d_%d", "file", rNamePre, size, deleteTime.Unix())
	return op.Mv(create.SourcePath, path.Join(clashDir, rName))
}

func getClashDir(realPath string) (string, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return "", err
	}
	for _, p := range partitions {
		if p.Mountpoint == "/" {
			continue
		}
		if strings.HasPrefix(realPath, p.Mountpoint) {
			clashDir := path.Join(p.Mountpoint, ".LinuxOnM_clash")
			if err = createClashDir(path.Join(p.Mountpoint, ".LinuxOnM_clash")); err != nil {
				return "", err
			}
			return clashDir, nil
		}
	}
	return constant.RecycleBinDir, createClashDir(constant.RecycleBinDir)
}

func createClashDir(clashDir string) error {
	op := files.NewFileOp()
	if !op.Stat(clashDir) {
		if err := op.CreateDir(clashDir, 0755); err != nil {
			return err
		}
	}
	return nil
}
