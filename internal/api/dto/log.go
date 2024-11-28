package dto

import "time"

type SearchLoginLogWithPage struct {
	PageInfo
	IP     string `json:"ip"`
	Status string `json:"status"`
}

type LoginLog struct {
	ID        uint      `json:"id"`
	IP        string    `json:"ip"`
	Address   string    `json:"address"`
	Agent     string    `json:"agent"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type SearchOpLogWithPage struct {
	PageInfo
	Source    string `json:"source"`
	Status    string `json:"status"`
	Operation string `json:"operation"`
}
