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

type OperationLog struct {
	ID     uint   `json:"id"`
	Source string `json:"source"`

	IP        string `json:"ip"`
	Path      string `json:"path"`
	Method    string `json:"method"`
	UserAgent string `json:"userAgent"`

	Latency time.Duration `json:"latency"`
	Status  string        `json:"status"`
	Message string        `json:"message"`

	DetailZH  string    `json:"detailZH"`
	DetailEN  string    `json:"detailEN"`
	CreatedAt time.Time `json:"createdAt"`
}

type SearchSSHLog struct {
	PageInfo
	Info   string `json:"info"`
	Status string `json:"Status" validate:"required,oneof=Success Failed All"`
}

type SSHLog struct {
	Logs            []SSHHistory `json:"logs"`
	TotalCount      int          `json:"totalCount"`
	SuccessfulCount int          `json:"successfulCount"`
	FailedCount     int          `json:"failedCount"`
}

type SSHHistory struct {
	Date     time.Time `json:"date"`
	DateStr  string    `json:"dateStr"`
	Area     string    `json:"area"`
	User     string    `json:"user"`
	AuthMode string    `json:"authMode"`
	Address  string    `json:"address"`
	Port     string    `json:"port"`
	Status   string    `json:"status"`
	Message  string    `json:"message"`
}
