package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"LinuxOnM/internal/utils/encrypt"
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

// SearchHost
// @Tags Host
// @Summary Page host
// @Description This function is used to retrieve a paginated list of hosts. It first validates and binds the incoming JSON request data of type dto.SearchHostWithPage.
//
//	The dto.SearchHostWithPage structure contains information such as the page number, page size, search criteria for hosts (in the 'Info' field), and the group ID to filter hosts by.
//	Then, it calls the SearchWithPage function in the HostService. In the HostService.SearchWithPage function, it queries the hostRepo to fetch the hosts based on the provided page number,
//	page size, and the search and group ID filters. If there's an error during this database query, it immediately returns with that error.
//	For each retrieved host, it copies the relevant data into a dto.HostInfo structure. It also fetches the corresponding group information for each host using the group ID and sets the
//	'GroupBelong' field in the dto.HostInfo structure to the group's name.
//	Additionally, if the 'RememberPassword' field of the host is false, the function clears the sensitive password-related fields (Password, PrivateKey, PassPhrase) in the dto.HostInfo structure.
//	However, if 'RememberPassword' is true and the respective password, private key, or pass phrase fields of the host have non-zero lengths, the function decrypts these values using appropriate
//	encryption functions before setting them in the dto.HostInfo structure. In case of any errors during the decryption process, it returns with that error.
//	Finally, it returns the total number of hosts that match the search criteria (regardless of pagination) and the list of hosts in the dto.HostInfo format.
//	Back in this route handling function, if the HostService.SearchWithPage call is successful, it packages the returned host list and total count into a dto.PageResult structure and sends it back
//	as a successful response. In case of any errors during the validation, binding, or the host retrieval and processing steps, appropriate error handling is performed and an error response is sent.
//
// @Accept json
// @Param request body dto.SearchHostWithPage true "request"
// @Success 200 {array} dto.HostTree
// @Security ApiKeyAuth
// @Router /host/search [post]
func (b *BaseApi) SearchHost(c *gin.Context) {
	var req dto.SearchHostWithPage
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	total, list, err := hostService.SearchWithPage(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, dto.PageResult{
		Items: list,
		Total: total,
	})
}

// UpdateHost
// @Tags Host
// @Summary Update host
// @Description This function is used to update host information. It first validates and binds the incoming JSON request data of type dto.HostOperate.
//
//	Then, it encrypts the password or private key (if provided) based on the authentication mode. If the authentication mode is 'password' and a password is present,
//	it calls the hostService.EncryptHost function to encrypt the password and clears the private key and pass phrase fields. If the authentication mode is 'key' and a
//	private key is present, it encrypts the private key and, if a pass phrase exists, encrypts it as well, and clears the password field.
//	Next, it constructs a map (upMap) containing the updated host information. The map includes fields such as name, group ID, address, port, user, authentication mode,
//	remember password flag, and description. The password, private key, and pass phrase fields are set according to the authentication mode.
//	Finally, it calls the hostService.Update function with the host ID and the update map. If the update is successful, it returns a success response with no data.
//	In case of any errors during the validation, binding, encryption, or update process, appropriate error handling is performed and an error response is sent.
//
// @Accept json
// @Param request body dto.HostOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/update [post]
// @x-panel-log {"bodyKeys":["name","addr"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新主机信息 [name][addr]","formatEN":"update host [name][addr]"}
func (b *BaseApi) UpdateHost(c *gin.Context) {
	var req dto.HostOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	var err error
	if len(req.Password) != 0 && req.AuthMode == "password" {
		req.Password, err = hostService.EncryptHost(req.Password)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
			return
		}
		req.PrivateKey = ""
		req.PassPhrase = ""
	}
	if len(req.PrivateKey) != 0 && req.AuthMode == "key" {
		req.PrivateKey, err = hostService.EncryptHost(req.PrivateKey)
		if err != nil {
			helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
			return
		}
		if len(req.PassPhrase) != 0 {
			req.PassPhrase, err = encrypt.StringEncrypt(req.PassPhrase)
			if err != nil {
				helper.ErrorWithDetail(c, constant.CodeErrBadRequest, constant.ErrTypeInvalidParams, err)
				return
			}
		}
		req.Password = ""
	}

	upMap := make(map[string]interface{})
	upMap["name"] = req.Name
	upMap["group_id"] = req.GroupID
	upMap["addr"] = req.Addr
	upMap["port"] = req.Port
	upMap["user"] = req.User
	upMap["auth_mode"] = req.AuthMode
	upMap["remember_password"] = req.RememberPassword
	if req.AuthMode == "password" {
		upMap["password"] = req.Password
		upMap["private_key"] = ""
		upMap["pass_phrase"] = ""
	} else {
		upMap["password"] = ""
		upMap["private_key"] = req.PrivateKey
		upMap["pass_phrase"] = req.PassPhrase
	}
	upMap["description"] = req.Description
	if err := hostService.Update(req.ID, upMap); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
