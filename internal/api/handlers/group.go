package handlers

import (
	"LinuxOnM/internal/api/dto"
	"LinuxOnM/internal/api/handlers/helper"
	"LinuxOnM/internal/constant"
	"github.com/gin-gonic/gin"
)

// ListGroup
// @Tags System Group
// @Summary List groups
// @Description This function is used to query and retrieve a list of system groups. It first validates and binds the incoming JSON request data of type dto.GroupSearch.
//
//	The dto.GroupSearch structure likely contains specific search criteria, such as the type of the group, which is used to filter the groups in the database.
//	Then, it calls the List function in the GroupService. In the GroupService.List function, it fetches the list of groups from the groupRepo based on the provided
//	search criteria (using the 'req.Type' value to filter by group type) and orders the results first by 'is_default' in descending order (so default groups come first)
//	and then by 'created_at' in descending order (newer groups first). If no groups are found during the retrieval process, it returns a specific error indicating that
//	the records were not found.
//	For each retrieved group, it copies the relevant data into a dto.GroupInfo structure and appends it to a slice. Finally, it returns an array of dto.GroupInfo
//	containing the information of the queried groups if the operation is successful. In case of any errors during the validation, binding, or the group retrieval and
//	transformation process, appropriate error handling is performed and an error response is sent.
//
// @Accept json
// @Param request body dto.GroupSearch true "request"
// @Success 200 {array} dto.GroupInfo
// @Security ApiKeyAuth
// @Router /group/search [post]
func (b *BaseApi) ListGroup(c *gin.Context) {
	var req dto.GroupSearch
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	list, err := groupService.List(req)
	if err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}

	helper.SuccessWithData(c, list)
}

// CreateGroup
// @Tags System Group
// @Summary Create group
// @Description This function serves the purpose of creating a new system group. It first validates and binds the incoming JSON request data of type dto.GroupCreate.
//
//	The dto.GroupCreate structure likely contains essential details required to create a group, such as the name and type of the group. These details are crucial for accurately creating the group record in the system.
//	Once the request data is successfully validated and bound, it proceeds to call the groupService.Create function, passing the validated request data (req) as an argument.
//	The groupService.Create function is responsible for performing the actual operations to create the group in the underlying database or relevant data storage. It might involve tasks like checking for uniqueness of the group name (if applicable), inserting the new group record with the provided details, and handling any associated business logic related to group creation.
//	If the group creation process within the groupService.Create function is executed without any errors, a success response with no additional data is returned to indicate that the group has been successfully created.
//	However, if any errors occur during the validation and binding of the request data or during the actual group creation process in the groupService.Create function, appropriate error handling is performed. In case of an error, the helper.ErrorWithDetail function is called to send back an error response. This error response includes detailed error information such as an error code (constant.CodeErrInternalServer) and an error type (constant.ErrTypeInternalServer), along with the specific error message that occurred during the process.
//
// @Accept json
// @Param request body dto.GroupCreate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /group [post]
// @x-panel-log {"bodyKeys":["name","type"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"创建组 [name][type]","formatEN":"create group [name][type]"}
func (b *BaseApi) CreateGroup(c *gin.Context) {
	var req dto.GroupCreate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := groupService.Create(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
