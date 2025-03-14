package models

type Command struct {
	BaseModel
	Name    string `gorm:"type:varchar(64);unique;not null" json:"name"`
	GroupID uint   `gorm:"type:decimal" json:"group_id"`
	Command string `gorm:"type:varchar(256);not null" json:"command"`
}
