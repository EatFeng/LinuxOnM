package services

import (
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/api/dto/response"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/cmd"
	"LinuxOnM/internal/utils/files"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type FileService struct{}

type IFileService interface {
	ReadLogByLine(req request.FileReadByLineReq) (*response.FileLineContent, error)
	GetFileList(op request.FileOption) (response.FileInfo, error)
}

func NewIFileService() IFileService {
	return &FileService{}
}

func (f *FileService) ReadLogByLine(req request.FileReadByLineReq) (*response.FileLineContent, error) {
	logFilePath := ""
	switch req.Type {
	case constant.TypeSystem:
		fileName := ""
		if len(req.Name) == 0 || req.Name == time.Now().Format("2006-01-02") {
			fileName = "LinuxOnM.log"
		} else {
			fileName = "LinuxOnM-" + req.Name + ".log"
		}
		logFilePath = path.Join(global.CONF.System.DataDir, "log", fileName)
		if _, err := os.Stat(logFilePath); err != nil {
			fileGzPath := path.Join(global.CONF.System.DataDir, "log", fileName+".gz")
			if _, err := os.Stat(fileGzPath); err != nil {
				return nil, buserr.New("ErrHttpReqNotFound")
			}
			if err := handleGunzip(fileGzPath); err != nil {
				return nil, fmt.Errorf("handle ungzip file %s failed, err: %v", fileGzPath, err)
			}
		}
	}

	lines, isEndOfFile, total, err := files.ReadFileByLine(logFilePath, req.Page, req.PageSize, req.Latest)
	if err != nil {
		return nil, err
	}
	if req.Latest && req.Page == 1 && len(lines) < 1000 && total > 1 {
		preLines, _, _, err := files.ReadFileByLine(logFilePath, total-1, req.PageSize, false)
		if err != nil {
			return nil, err
		}
		lines = append(preLines, lines...)
	}

	res := &response.FileLineContent{
		Content: strings.Join(lines, "\n"),
		End:     isEndOfFile,
		Path:    logFilePath,
		Total:   total,
		Lines:   lines,
	}
	return res, nil
}

func handleGunzip(path string) error {
	if _, err := cmd.Execf("gunzip %s", path); err != nil {
		return err
	}
	return nil
}

func (f *FileService) GetFileList(op request.FileOption) (response.FileInfo, error) {
	var fileInfo response.FileInfo
	data, err := os.Stat(op.Path)
	if err != nil && os.IsNotExist(err) {
		return fileInfo, nil
	}
	if !data.IsDir() {
		op.FileOption.Path = filepath.Dir(op.FileOption.Path)
	}
	info, err := files.NewFileInfo(op.FileOption)
	if err != nil {
		return fileInfo, err
	}
	fileInfo.FileInfo = *info
	return fileInfo, nil
}
