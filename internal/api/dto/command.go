package dto

type CommandInfo struct {
	ID          uint   `json:"id"`
	GroupID     uint   `json:"group_id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	GroupBelong string `json:"groupBelong"`
}

type CommandOperate struct {
	ID          uint   `json:"id"`
	GroupID     uint   `json:"group_id"`
	GroupBelong string `json:"groupBelong"`
	Name        string `json:"name" validate:"required"`
	Command     string `json:"command" validate:"required"`
}
