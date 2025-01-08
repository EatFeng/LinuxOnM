package request

type ProcessReq struct {
	PID int32 `json:"PID"  validate:"required"`
}

type ProcessRequest struct {
	Name string `json:"name" validate:"required"`
}

type ProcessCreate struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}
