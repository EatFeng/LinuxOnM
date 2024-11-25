package models

type LoginLog struct {
	BaseModel
	IP      string `gorm:"type:varchar(64)" json:"ip"`
	Address string `gorm:"type:varchar(64)" json:"address"`
	Agent   string `gorm:"type:varchar(64)" json:"agent"`
	Status  string `gorm:"type:varchar(64)" json:"status"`
	Message string `gorm:"type:text" json:"message"`
}
