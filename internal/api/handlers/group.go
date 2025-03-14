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

// UpdateGroup
// @Tags System Group
// @Summary Update group
// @Description This function is responsible for handling the update operation of system groups. It first validates and binds the incoming JSON request data of type dto.GroupUpdate.
//
//	The dto.GroupUpdate structure likely contains fields that are relevant for modifying the properties of an existing group, such as the 'name' and 'type' fields which are specified in the @x-panel-log annotation as key fields. These fields hold the updated values that the user intends to apply to the corresponding group.
//	Once the request data is successfully bound to the 'req' variable of type dto.GroupUpdate, the function proceeds to call the groupService.Update function, passing the 'req' object as an argument.
//	The groupService.Update function is tasked with performing the actual update operations in the underlying database or relevant data storage. This may involve checking the validity of the new values (for example, ensuring that the group name adheres to any naming conventions or uniqueness constraints if applicable), querying the database to locate the group record to be updated using its unique identifier (which might be implicitly associated with the request data), and then updating the corresponding fields in the group record with the new values provided in the 'req' object.
//	If the update process within the groupService.Update function is completed without encountering any errors, a success response with no additional data is returned to indicate that the group has been successfully updated.
//	However, if any issues arise during the validation and binding of the request data (such as incorrect JSON formatting, or the data not conforming to the validation rules defined for the dto.GroupUpdate structure) or during the actual update process in the groupService.Update function (like database connection failures, or violations of data integrity constraints), appropriate error handling is carried out. In such cases, the helper.ErrorWithDetail function is called to send back an error response. This error response includes detailed error information such as an error code (constant.CodeErrInternalServer) and an error type (constant.ErrTypeInternalServer), along with the specific error message that occurred during the process.
//
// @Accept json
// @Param request body dto.GroupUpdate true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /group/update [post]
// @x-panel-log {"bodyKeys":["name","type"],"paramKeys":[],"BeforeFunctions":[],"formatZH":"更新组 [name][type]","formatEN":"update group [name][type]"}
func (b *BaseApi) UpdateGroup(c *gin.Context) {
	var req dto.GroupUpdate
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := groupService.Update(req); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}

// DeleteGroup
// @Tags System Group
// @Summary Delete group
// @Description This function is specifically designed to handle the deletion of system groups. It commences by validating and binding the incoming JSON request data of type dto.OperateByID.
//
//	The dto.OperateByID structure is likely crafted to hold the essential identifier information required to pinpoint the specific system group that needs to be deleted. In this case, it presumably contains an 'ID' field which serves as the unique identifier for the target group within the system.
//	Once the request data has been successfully validated and bound to the 'req' variable of type dto.OperateByID, the function proceeds to call the groupService.Delete function, passing the group's ID (req.ID) as an argument.
//	The groupService.Delete function is responsible for executing the actual deletion operations in the underlying database or relevant data storage. This entails several important steps, such as first verifying whether the group can be safely deleted. This might involve checking if the group is not associated with any other critical entities or operations that could be disrupted by its removal (for example, ensuring there are no dependencies on this group from other parts of the system). After confirming that the deletion is feasible, it proceeds to perform the necessary database operations to physically remove the group record from the storage.
//	If the deletion process within the groupService.Delete function is carried out without any errors, a success response with no additional data is returned to signify that the system group has been successfully deleted.
//	However, if any issues arise during the validation and binding of the request data (for instance, if the JSON format of the incoming request is incorrect or the 'ID' value doesn't meet the validation requirements specified for the dto.OperateByID structure) or during the actual deletion process in the groupService.Delete function (like encountering database connection problems, violating deletion constraints, or issues with data integrity), appropriate error handling is implemented. In such situations, the helper.ErrorWithDetail function is invoked to send back an error response. This error response encompasses detailed error information including an error code (constant.CodeErrInternalServer) and an error type (constant.ErrTypeInternalServer), along with the specific error message that occurred during the process.
//
// @Accept json
// @Param request body dto.OperateByID true "request"
// @Success 200
// @Security ApiKeyAuth
// @Router /group/del [post]
// @x-panel-log {"bodyKeys":["id"],"paramKeys":[],"BeforeFunctions":[{"input_column":"id","input_value":"id","isList":false,"db":"groups","output_column":"name","output_value":"name"},{"input_column":"id","input_value":"id","isList":false,"db":"groups","output_column":"type","output_value":"type"}],"formatZH":"删除组 [type][name]","formatEN":"delete group [type][name]"}
func (b *BaseApi) DeleteGroup(c *gin.Context) {
	var req dto.OperateByID
	if err := helper.CheckBindAndValidate(c, &req); err != nil {
		return
	}

	if err := groupService.Delete(req.ID); err != nil {
		helper.ErrorWithDetail(c, constant.CodeErrInternalServer, constant.ErrTypeInternalServer, err)
		return
	}
	helper.SuccessWithData(c, nil)
}
