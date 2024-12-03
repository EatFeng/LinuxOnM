package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// TestByID
// @Tags Host
// @Summary Test host conn by host id
// @Description 测试主机连接
// @Accept json
// @Param id path integer true "request"
// @Success 200 {boolean} connStatus
// @Security ApiKeyAuth
// @Router /host/test/byid/:id [post]
func (b *BaseApi) TestByID(c *gin.Context) {
	id, err := helper.GetParamID(c)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
		return
	}

	connStatus := hostService.TestLocalConn(id)
	helper.SuccessWithData(c, connStatus)
}

// TestByInfo
// @Tags Host
// @Summary Test host connection by provided connection information
// @Description This function is used to test the SSH connection to a host based on the connection information provided in the request body.
//
//	It first validates and binds the incoming JSON request data to the dto.HostConnTest structure. Then, it calls the TestByInfo
//	function in the HostService to perform the actual connection test. If the connection test is successful, it returns a 200 status
//	code along with the connection status (true). In case of any errors during validation, binding, or the connection test itself,
//	appropriate error handling is performed.
//
// @Accept json
// @Param request body dto.HostConnTest true "request"
// @Success 200 "Returns true if the SSH connection to the host is successfully established, false otherwise."
// @Security ApiKeyAuth
// @Router /hosts/test/byinfo [post]
func (b *BaseApi) TestByInfo(c *gin.Context) {
	var req dto.HostConnTest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	connStatus := hostService.TestByInfo(req)
	helper.SuccessWithData(c, connStatus)
}
