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
// @Router /file/search [post]
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
// @Router /file [post]
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
// @Router /file/del [post]
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
// @Router /file/upload [post]
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

// @Tags File
// @Summary Load file content
// @Accept json
// @Param request body request.FileContentReq true "request"
// @Success 200 {object} response.FileInfo
// @Security ApiKeyAuth
// @Router /file/content [post]
// @x-panel-log {"bodyKeys":["path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"获取文件内容 [path]","formatEN":"Load file content [path]"}
func (b *BaseApi) GetContent(c *gin.Context) {
	var req request.FileContentReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	info, err := fileService.GetContent(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, info)
}

// @Tags File
// @Summary Batch change file mode and owner
// @Accept json
// @Param request body request.FileRoleReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/batch/role [post]
// @x-panel-log {"bodyKeys":["paths","mode","user","group"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"批量修改文件权限和用户/组 [paths] => [mode]/[user]/[group]","formatEN":"Batch change file mode and owner [paths] => [mode]/[user]/[group]"}
func (b *BaseApi) BatchChangeModeAndOwner(c *gin.Context) {
	var req request.FileRoleReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := fileService.BatchChangeModeAndOwner(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
	}
	helper.SuccessWithOutData(c)
}

// @Tags File
// @Summary Check file exist
// @Accept json
// @Param request body request.FilePathCheck true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/check [post]
func (b *BaseApi) CheckFile(c *gin.Context) {
	var req request.FilePathCheck
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if _, err := os.Stat(req.Path); err != nil {
		helper.SuccessWithData(c, false)
		return
	}
	helper.SuccessWithData(c, true)
}

// @Tags File
// @Summary Move file
// @Accept json
// @Param request body request.FileMove true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/move [post]
// @x-panel-log {"bodyKeys":["oldPaths","newPath"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"移动文件 [oldPaths] => [newPath]","formatEN":"Move [oldPaths] => [newPath]"}
func (b *BaseApi) MoveFile(c *gin.Context) {
	var req request.FileMove
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := fileService.MvFile(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags File
// @Summary Change file name
// @Accept json
// @Param request body request.FileRename true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/rename [post]
// @x-panel-log {"bodyKeys":["oldName","newName"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"重命名 [oldName] => [newName]","formatEN":"Rename [oldName] => [newName]"}
func (b *BaseApi) ChangeFileName(c *gin.Context) {
	var req request.FileRename
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	if err := fileService.ChangeName(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags File
// @Summary Load files tree
// @Accept json
// @Param request body request.FileOption true "request"
// @Success 200 {array} response.FileTree
// @Security ApiKeyAuth
// @Router /file/tree [post]
func (b *BaseApi) GetFileTree(c *gin.Context) {
	var req request.FileOption
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	tree, err := fileService.GetFileTree(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, tree)
}

// @Tags File
// @Summary Compress file
// @Accept json
// @Param request body request.FileCompress true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/compress [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"压缩文件 [name]","formatEN":"Compress file [name]"}
func (b *BaseApi) CompressFile(c *gin.Context) {
	var req request.FileCompress
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	err := fileService.Compress(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags File
// @Summary Decompress file
// @Accept json
// @Param request body request.FileDeCompress true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /file/decompress [post]
// @x-panel-log {"bodyKeys":["path"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"解压 [path]","formatEN":"Decompress file [path]"}
func (b *BaseApi) DeCompressFile(c *gin.Context) {
	var req request.FileDeCompress
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	err := fileService.DeCompress(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
