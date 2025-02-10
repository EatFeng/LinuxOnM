package dto

type PageInfo struct {
	Page     int `json:"page" validate:"required,number"`
	PageSize int `json:"pageSize" validate:"required,number"`
}

type BatchDeleteReq struct {
	Ids []uint `json:"ids" validate:"required"`
}

type OperateByID struct {
	ID uint `json:"id" validate:"required"`
}

type SearchWithPage struct {
	PageInfo
	Info string `json:"info"`
}
