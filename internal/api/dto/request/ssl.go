package request

type WebsiteCACreate struct {
	CommonName       string `json:"commonName" validate:"required"`
	Country          string `json:"country" validate:"required"`
	Organization     string `json:"organization" validate:"required"`
	OrganizationUint string `json:"organizationUint"`
	Name             string `json:"name" validate:"required"`
	KeyType          string `json:"keyType" validate:"required,oneof=P256 P384 2048 3072 4096 8192"`
	Province         string `json:"province" `
	City             string `json:"city"`
}
