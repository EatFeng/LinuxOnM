package dto

type GroupSearch struct {
	Type string `json:"type" validate:"required"`
}

type GroupInfo struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"isDefault"`
}
