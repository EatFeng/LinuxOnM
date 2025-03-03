package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// ContainerStats
// @Tags Container
// @Summary Container stats
// @Description 容器监控信息
// @Param id path integer true "容器id"
// @Success 200 {object} dto.ContainerStats
// @Security ApiKeyAuth
// @Router /container/stats/:id [get]
func (b *BaseApi) ContainerStats(c *gin.Context) {
	containerID, ok := c.Params.Get("id")
	if !ok {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, errors.New("error container id in path"))
		return
	}

	result, err := containerService.ContainerStats(containerID)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, result)
}

// @Tags Container
// @Summary Create container
// @Description 创建容器
// @Accept json
// @Param request body dto.ContainerOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /container [post]
// @x-panel-log {"bodyKeys":["name","image"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建容器 [name][image]","formatEN":"create container [name][image]"}
func (b *BaseApi) ContainerCreate(c *gin.Context) {
	var req dto.ContainerOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerCreate(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Container
// @Summary Update container
// @Description 更新容器
// @Accept json
// @Param request body dto.ContainerOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /containers/update [post]
// @x-panel-log {"bodyKeys":["name","image"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新容器 [name][image]","formatEN":"update container [name][image]"}
func (b *BaseApi) ContainerUpdate(c *gin.Context) {
	var req dto.ContainerOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerUpdate(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Container
// @Summary Upgrade container
// @Description 更新容器镜像
// @Accept json
// @Param request body dto.ContainerUpgrade true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /container/upgrade [post]
// @x-panel-log {"bodyKeys":["name","image"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新容器镜像 [name][image]","formatEN":"upgrade container image [name][image]"}
func (b *BaseApi) ContainerUpgrade(c *gin.Context) {
	var req dto.ContainerUpgrade
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerUpgrade(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// ContainerInfo
// @Tags Container
// @Summary Load container info
// @Description 获取容器表单信息
// @Accept json
// @Param request body dto.OperationWithName true "request"
// @Success 200 {object} dto.ContainerOperate
// @Security ApiKeyAuth
// @Router /container/info [post]
func (b *BaseApi) ContainerInfo(c *gin.Context) {
	var req dto.OperationWithName
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	data, err := containerService.ContainerInfo(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, data)
}

// SearchContainer
// @Tags Container
// @Summary Page containers
// @Description 获取容器列表分页
// @Accept json
// @Param request body dto.PageContainer true "request"
// @Produce json
// @Success 200 {object} dto.PageResult
// @Security ApiKeyAuth
// @Router /container/search [post]
func (b *BaseApi) SearchContainer(c *gin.Context) {
	var req dto.PageContainer
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	total, list, err := containerService.Page(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// ListContainer
// @Tags Container
// @Summary List containers
// @Description 获取容器名称
// @Accept json
// @Produce json
// @Success 200
// @Security ApiKeyAuth
// @Router /container/list [post]
func (b *BaseApi) ListContainer(c *gin.Context) {
	list, err := containerService.List()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, list)
}

// ContainerListStats
// @Summary Load container stats
// @Description 获取容器列表资源占用
// @Success 200 {array} dto.ContainerListStats
// @Security ApiKeyAuth
// @Router /container/list/stats [get]
func (b *BaseApi) ContainerListStats(c *gin.Context) {
	data, err := containerService.ContainerListStats()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, data)
}

// @Tags Container
// @Summary Operate Container
// @Description 容器操作
// @Accept json
// @Param request body dto.ContainerOperation true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /containers/operate [post]
// @x-panel-log {"bodyKeys":["names","operation"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"容器 [names] 执行 [operation]","formatEN":"container [operation] [names]"}
func (b *BaseApi) ContainerOperation(c *gin.Context) {
	var req dto.ContainerOperation
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerOperation(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// ContainerLogs
// @Tags Container
// @Summary Container logs
// @Description 容器日志
// @Param container query string false "容器名称"
// @Param since query string false "时间筛选"
// @Param follow query string false "是否追踪"
// @Param tail query string false "显示行号"
// @Security ApiKeyAuth
// @Router /container/search/log [get]
func (b *BaseApi) ContainerLogs(c *gin.Context) {
	wsConn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.LOG.Errorf("gin context http handler failed, err: %v", err)
		return
	}
	defer wsConn.Close()

	container := c.Query("container")
	since := c.Query("since")
	follow := c.Query("follow") == "true"
	tail := c.Query("tail")

	if err := containerService.ContainerLogs(wsConn, "container", container, since, tail, follow); err != nil {
		_ = wsConn.WriteMessage(1, []byte(err.Error()))
		return
	}
}

// @Description 下载容器日志
// @Router /containers/download/log [post]
func (b *BaseApi) DownloadContainerLogs(c *gin.Context) {
	var req dto.ContainerLog
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}
	err := containerService.DownloadContainerLogs(req.ContainerType, req.Container, req.Since, strconv.Itoa(int(req.Tail)), c)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
	}
}

// @Tags Container
// @Summary Clean container log
// @Description 清理容器日志
// @Accept json
// @Param request body dto.OperationWithName true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /containers/clean/log [post]
// @x-panel-log {"bodyKeys":["name"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"清理容器 [name] 日志","formatEN":"clean container [name] logs"}
func (b *BaseApi) CleanContainerLog(c *gin.Context) {
	var req dto.OperationWithName
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerLogClean(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// LoadResourceLimit
// @Summary Load container limits
// @Description 获取容器资源限制
// @Success 200 {object} dto.ResourceLimit
// @Security ApiKeyAuth
// @Router /container/limit [get]
func (b *BaseApi) LoadResourceLimit(c *gin.Context) {
	data, err := containerService.LoadResourceLimit()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, data)
}

// @Tags Container
// @Summary Container inspect
// @Description 容器详情
// @Accept json
// @Param request body dto.InspectReq true "request"
// @Success 200 {string} result
// @Security ApiKeyAuth
// @Router /container/inspect [post]
func (b *BaseApi) Inspect(c *gin.Context) {
	var req dto.InspectReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	result, err := containerService.Inspect(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, result)
}

// @Tags Container
// @Summary Rename Container
// @Description 容器重命名
// @Accept json
// @Param request body dto.ContainerRename true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /container/rename [post]
// @x-panel-log {"bodyKeys":["name","newName"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"容器重命名 [name] => [newName]","formatEN":"rename container [name] => [newName]"}
func (b *BaseApi) ContainerRename(c *gin.Context) {
	var req dto.ContainerRename
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerRename(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Container
// @Summary Commit Container
// @Description 容器提交生成新镜像
// @Accept json
// @Param request body dto.ContainerCommit true "request"
// @Success 200
// @Router /container/commit [post]
func (b *BaseApi) ContainerCommit(c *gin.Context) {
	var req dto.ContainerCommit
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := containerService.ContainerCommit(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// @Tags Container
// @Summary Clean container
// @Description 容器清理
// @Accept json
// @Param request body dto.ContainerPrune true "request"
// @Success 200 {object} dto.ContainerPruneReport
// @Security ApiKeyAuth
// @Router /container/prune [post]
// @x-panel-log {"bodyKeys":["pruneType"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"清理容器 [pruneType]","formatEN":"clean container [pruneType]"}
func (b *BaseApi) ContainerPrune(c *gin.Context) {
	var req dto.ContainerPrune
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	report, err := containerService.Prune(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, report)
}

// ListNetwork
// @Tags Container Network
// @Summary List networks
// @Description 获取容器网络列表
// @Accept json
// @Produce json
// @Success 200 {array} dto.Options
// @Security ApiKeyAuth
// @Router /container/network [get]
func (b *BaseApi) ListNetwork(c *gin.Context) {
	list, err := containerService.ListNetwork()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, list)
}

// ListVolume
// @Tags Container Volume
// @Summary List volumes
// @Description 获取容器存储卷列表
// @Accept json
// @Produce json
// @Success 200 {array} dto.Options
// @Security ApiKeyAuth
// @Router /container/volume [get]
func (b *BaseApi) ListVolume(c *gin.Context) {
	list, err := containerService.ListVolume()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, list)
}
