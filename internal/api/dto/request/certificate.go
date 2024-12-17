package request

type CertificateObtain struct {
	ID          uint   `json:"id" validate:"required"`
	Domains     string `json:"domains" validate:"required"`
	KeyType     string `json:"keyType" validate:"required,oneof=P256 P384 2048 3072 4096 8192"`
	Time        int    `json:"time" validate:"required"`
	Unit        string `json:"unit" validate:"required"`
	PushDir     bool   `json:"pushDir"`
	Dir         string `json:"dir"`
	AutoRenew   bool   `json:"autoRenew"`
	Renew       bool   `json:"renew"`
	SSLID       uint   `json:"sslID"`
	Description string `json:"description"`
	ExecShell   bool   `json:"execShell"`
	Shell       string `json:"shell"`
}
