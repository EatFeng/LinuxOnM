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

// UpdateHostGroup
// @Tags Host
// @Summary Update host group
// @Description This function is designed to handle the operation of changing the group to which a host belongs, which is often referred to as "switching the host's group".
//
//	It first validates and binds the incoming JSON request data of type dto.ChangeHostGroup. The dto.ChangeHostGroup structure likely contains essential information such as the ID of the host
//	whose group needs to be changed and the ID of the target group it will be switched to.
//	Next, it constructs a map named 'upMap' which is used to hold the update information. In this case, only the 'group_id' key is set with the value from req.GroupID, indicating the new group ID that the host will be assigned to.
//	Then, it calls the hostService.Update function with the host's ID (req.ID) and the 'upMap'. This function in the hostService is responsible for performing the actual update operation in the database to change the host's group.
//	If the update process is successful, a success response with no additional data is returned. However, if any errors occur during the validation and binding of the request data or the database update operation,
//	appropriate error handling is performed. Specifically, if there's an error during the update, the helper.ErrorWithDetail function is called to send back an error response with detailed error information, including an error code
//	(constant.CodeErrInternalServer) and an error type (constant.ErrTypeInternalServer), along with the specific error message.
//
// @Accept json
// @Param request body dto.ChangeHostGroup true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/update/group [post]
// @x-panel-log {"bodyKeys":["id","group"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"hosts","output_column":"addr","output_value":"addr"}],"formatZH":"切换主机[addr]分组 => [group]","formatEN":"change host [addr] group => [group]"}
func (b *BaseApi) UpdateHostGroup(c *gin.Context) {
	var req dto.ChangeHostGroup
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	upMap := make(map[string]interface{})
	upMap["group_id"] = req.GroupID
	if err := hostService.Update(req.ID, upMap); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// DeleteHost
// @Tags Host
// @Summary Delete host
// @Description This function is designed to handle the deletion of hosts. It operates in a way that first validates and binds the incoming JSON request data of type dto.BatchDeleteReq.
//
//	The dto.BatchDeleteReq structure is expected to contain a specific field (as per its definition) which is likely used to identify the hosts that need to be deleted. In this case, it probably has an 'Ids' field that holds the identifiers (such as unique IDs) of the hosts targeted for deletion.
//	Once the request data has been successfully validated and bound to the 'req' variable of type dto.BatchDeleteReq, the function proceeds to call the hostService.Delete function, passing the list of host IDs (req.Ids) as an argument.
//	The hostService.Delete function is responsible for executing the actual deletion operations in the underlying database or relevant data storage. This might involve performing necessary checks, such as ensuring that the hosts are eligible for deletion (for example, checking if they are not referenced by other critical components or operations), and then physically removing the corresponding host records from the database.
//	If the deletion process within the hostService.Delete function is completed without encountering any errors, a success response with no additional data is returned to indicate that the hosts have been successfully deleted.
//	However, if any issues arise during the validation and binding of the request data (for instance, if the format of the incoming JSON is incorrect or the data doesn't meet the validation requirements specified for the dto.BatchDeleteReq structure) or during the actual deletion process in the hostService.Delete function (like database connection problems or violations of deletion constraints), appropriate error handling is carried out. In such cases, the helper.ErrorWithDetail function is called to send back an error response. This error response includes detailed error information such as an error code (constant.CodeErrInternalServer) and an error type (constant.ErrTypeInternalServer), along with the specific error message that occurred during the process.
//
// @Accept json
// @Param request body dto.BatchDeleteReq true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/del [post]
// @x-panel-log {"bodyKeys":["ids"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"ids","isList":true,"db":"hosts","output_column":"addr","output_value":"addrs"}],"formatZH":"删除主机 [addrs]","formatEN":"delete host [addrs]"}
func (b *BaseApi) DeleteHost(c *gin.Context) {
	var req dto.BatchDeleteReq
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := hostService.Delete(req.Ids); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// ListCommand
// @Tags Command
// @Summary List commands
// @Description This function is used to retrieve the list of quick commands. It calls the commandService.List function which is responsible for querying and fetching the relevant command data from the underlying data source, such as a database. If the query operation in the commandService.List function is successful, it returns the list of commands which is then sent back as a successful response with a status code of 200 and the data in the format of dto.CommandInfo. In case of any errors during the query process, such as database connection issues or errors in data retrieval, the helper.ErrorWithDetail function is called to send back an error response with a specific error code (constant.CodeErrInternalServer) and error type (constant.ErrTypeInternalServer), along with the detailed error message.
// @Success 200 {object} dto.CommandInfo
// @Security ApiKeyAuth
// @Router /host/command [get]
func (b *BaseApi) ListCommand(c *gin.Context) {
	list, err := commandService.List()
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, list)
}

// CreateCommand
// @Tags Command
// @Summary Create command
// @Description This function is designed to create a new quick command. It first validates and binds the incoming JSON request data of type dto.CommandOperate. The dto.CommandOperate structure likely contains essential fields for creating a command, such as the name and the actual command text. After successful validation and binding, it calls the commandService.Create function, passing the validated request data (req) as an argument. The commandService.Create function is tasked with performing the actual creation operations in the underlying data source, which may involve inserting the new command record with the provided details and handling any associated business logic or data integrity checks. If the creation process is successful, a success response with no additional data is returned. In case of any errors during the validation and binding of the request data or during the actual creation process in the commandService.Create function, the helper.ErrorWithDetail function is called to send back an error response with the appropriate error code and type, along with the detailed error message.
// @Accept json
// @Param request body dto.CommandOperate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /host/command [post]
// @x-panel-log {"bodyKeys":["name","command"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建快捷命令 [name][command]","formatEN":"create quick command [name][command]"}
func (b *BaseApi) CreateCommand(c *gin.Context) {
	var req dto.CommandOperate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := commandService.Create(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
