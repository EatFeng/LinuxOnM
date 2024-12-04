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
