package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/global"

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
