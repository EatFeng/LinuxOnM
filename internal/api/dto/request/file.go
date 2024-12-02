package request

type FileReadByLineReq struct {
	Page     int    `json:"page" validate:"required"`
	PageSize int    `json:"pageSize" validate:"required"`
	Type     string `json:"type" validate:"required"`
	ID       uint   `json:"ID"`
	Name     string `json:"name"`
	Latest   bool   `json:"latest"`
}
