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
// @Router /host/test/byinfo [post]
func (b *BaseApi) TestByInfo(c *gin.Context) {
	var req dto.HostConnTest
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	connStatus := hostService.TestByInfo(req)
	helper.SuccessWithData(c, connStatus)
}

// CreateHost
// @Tags Host
// @Summary Create host
// @Description This function is used to create or update a host. It first validates and binds the incoming JSON request data of type dto.HostOperate.
//
//	Then it calls the Create function in the HostService. In the HostService.Create function, it encrypts the password or private key if needed
//	based on the authentication mode. It also determines the host group, checks if a host with the same address (and additional criteria if not 127.0.0.1)
//	already exists. If it exists, it updates the existing host record; otherwise, it creates a new one. Finally, it returns the created or updated host
//	information in the dto.HostInfo type. If any errors occur during the process, appropriate error handling is performed and an error response is sent.
//
// @Accept json
// @Param request body dto.HostOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host [post]
// @x-panel-log {"bodyKeys":["name","addr"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建主机 [name][addr]","formatEN":"create host [name][addr]"}
func (b *BaseApi) CreateHost(c *gin.Context) {
	var req dto.HostOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	host, err := hostService.Create(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, host)
}

// HostTree
// @Tags Host
// @Summary Load host tree
// @Description This function is used to load the host information in a tree-like structure. It first validates and binds the incoming JSON request data of type dto.SearchForTree.
//
//	The dto.SearchForTree structure likely contains specific search criteria to filter the hosts. Then, it calls the SearchForTree function in the HostService.
//	In the HostService.SearchForTree function, it retrieves the list of hosts from the database based on the provided search information and also fetches the list of host groups.
//	It then constructs a hierarchical tree structure where each host group is a node and the hosts belonging to that group are its children. The details of each host node include its ID and a label that combines information like host name (if exists), user, address, and port.
//	If the construction of the tree structure is successful and no errors occur during the database operations and data processing, it returns an array of dto.HostTree containing the complete host tree information.
//	In case of any errors during the process, appropriate error handling is performed and an error response is sent.
//
// @Accept json
// @Param request body dto.SearchForTree true "request"
// @Success 200 {array} dto.HostTree
// @Security ApiKeyAuth
// @Router /host/tree [post]
func (b *BaseApi) HostTree(c *gin.Context) {
	var req dto.SearchForTree
	if err := helper.CheckBind(&req, c); err != nil {
		return
	}

	data, err := hostService.SearchForTree(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, data)
}
