package handlers

import (
	"LinuxOnM/internal/api/dto/request"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/buserr"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"LinuxOnM/internal/utils/files"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
)

// @Tags File
// @Summary List files
// @Accept json
// @Param request body request.FileOption true "request"
// @Success 200 {object} response.FileInfo
// @Security ApiKeyAuth
// @Router /files/search [post]
func (b *BaseApi) ListFiles(c *gin.Context) {
	var req request.FileOption
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	files, err := fileService.GetFileList(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, files)
}

// @Tags File
// @Summary Read file by Line
// @Description 按行读取日志文件
// @Param request body request.FileReadByLineReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/read [post]
func (b *BaseApi) ReadFileByLine(c *gin.Context) {
	var req request.FileReadByLineReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	res, err := fileService.ReadLogByLine(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, res)
}

var wsUpgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Tags File
// @Summary Create file
// @Accept json
// @Param request body request.FileCreate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /files [post]
// @x-panel-log {"bodyKeys":["path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建文件/文件夹 [path]","formatEN":"Create dir or file [path]"}
func (b *BaseApi) CreateFile(c *gin.Context) {
	var req request.FileCreate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	err := fileService.Create(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags File
// @Summary Delete file
// @Accept json
// @Param request body request.FileDelete true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /files/del [post]
// @x-panel-log {"bodyKeys":["path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"删除文件/文件夹 [path]","formatEN":"Delete dir or file [path]"}
func (b *BaseApi) DeleteFile(c *gin.Context) {
	var req request.FileDelete
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	err := fileService.Delete(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags File
// @Summary Upload file
// @Param file formData file true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /files/upload [post]
// @x-panel-log {"bodyKeys":["path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"上传文件 [path]","formatEN":"Upload file [path]"}
func (b *BaseApi) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}
	uploadFiles := form.File["file"]
	paths := form.Value["path"]

	overwrite := true
	if ow, ok := form.Value["overwrite"]; ok {
		if len(ow) != 0 {
			parseBool, _ := strconv.ParseBool(ow[0])
			overwrite = parseBool
		}
	}

	if len(paths) == 0 || !strings.Contains(paths[0], "/") {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, errors.New("error paths in request"))
		return
	}
	dir := path.Dir(paths[0])

	_, err = os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		mode, err := files.GetParentMode(dir)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
			return
		}
		if err = os.MkdirAll(dir, mode); err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, fmt.Errorf("mkdir %s failed, err: %v", dir, err))
			return
		}
	}
	info, err := os.Stat(dir)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	mode := info.Mode()
	fileOp := files.NewFileOp()
	stat, ok := info.Sys().(*syscall.Stat_t)
	uid, gid := -1, -1
	if ok {
		uid, gid = int(stat.Uid), int(stat.Gid)
	}

	success := 0
	failures := make(buserr.MultiErr)
	for _, file := range uploadFiles {
		dstFilename := path.Join(paths[0], file.Filename)
		dstDir := path.Dir(dstFilename)
		if !fileOp.Stat(dstDir) {
			if err = fileOp.CreateDir(dstDir, mode); err != nil {
				e := fmt.Errorf("create dir [%s] failed, err: %v", path.Dir(dstFilename), err)
				failures[file.Filename] = e
				global.LOG.Error(e)
				continue
			}
			_ = os.Chown(dstDir, uid, gid)
		}
		tmpFilename := dstFilename + ".tmp"
		if err := c.SaveUploadedFile(file, tmpFilename); err != nil {
			_ = os.Remove(tmpFilename)
			e := fmt.Errorf("upload [%s] file failed, err: %v", file.Filename, err)
			failures[file.Filename] = e
			global.LOG.Error(e)
			continue
		}
		dstInfo, statErr := os.Stat(dstFilename)
		if overwrite {
			_ = os.Remove(dstFilename)
		}

		err = os.Rename(tmpFilename, dstFilename)
		if err != nil {
			_ = os.Remove(tmpFilename)
			e := fmt.Errorf("upload [%s] file failed, err: %v", file.Filename, err)
			failures[file.Filename] = e
			global.LOG.Error(e)
			continue
		}
		if statErr == nil {
			_ = os.Chmod(dstFilename, dstInfo.Mode())
		} else {
			_ = os.Chmod(dstFilename, mode)
		}
		if uid != -1 && gid != -1 {
			_ = os.Chown(dstFilename, uid, gid)
		}
		success++
	}
	if success == 0 {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, failures)
	} else {
		helper.SuccessWithMsg(c, fmt.Sprintf("%d files upload success", success))
	}
}
