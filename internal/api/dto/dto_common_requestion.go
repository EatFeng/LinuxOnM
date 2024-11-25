package dto

type PageInfo struct {
	Page     int `json:"page" validate:"required,number"`
	PageSize int `json:"pageSize" validate:"required,number"`
}
