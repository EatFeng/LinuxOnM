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

type SearchCommandWithPage struct {
	PageInfo
	OrderBy string `json:"orderBy" validate:"required,oneof=name command created_at"`
	Order   string `json:"order" validate:"required,oneof=null ascending descending"`
	GroupID uint   `json:"group_id"`
	Info    string `json:"info"`
	Name    string `json:"name"`
}

type CommandTree struct {
	ID       uint          `json:"id"`
	Label    string        `json:"label"`
	Children []CommandInfo `json:"children"`
}
