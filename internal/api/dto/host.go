package dto

type HostConnTest struct {
	Addr       string `json:"addr" validate:"required"`
	Port       uint   `json:"port" validate:"required,number,max=65535,min=1"`
	User       string `json:"user" validate:"required"`
	AuthMode   string `json:"authMode" validate:"oneof=password key"`
	Password   string `json:"password"`
	PrivateKey string `json:"privateKey"`
	PassPhrase string `json:"passPhrase"`
}
