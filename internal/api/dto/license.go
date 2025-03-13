package dto

type LicenseUploadRequest struct {
	LicenseFile []byte `json:"licenseFile" binding:"required"` // 上传的.lic文件
	Password    string `json:"password"`                       // AES解密密码（可选）
}

type LicenseUploadResponse struct {
	LicenseID    string `json:"licenseId"`
	ExpiryDate   string `json:"expiryDate"`
	HardwareHash string `json:"hardwareHash"`
	IsValid      bool   `json:"isValid"`
}

type SignedLicenseData struct {
	ExpiryDate   string `json:"expiry_date"` // 调整字段顺序
	HardwareHash string `json:"hardware_hash"`
	IssuedAt     int    `json:"issued_at"`
	LicenseID    string `json:"license_id"` // 按字母顺序最后
}

type LicenseStatusResponse struct {
	LicenseID      string `json:"license_id"`
	ExpiryDate     string `json:"expiry_date"`
	HardwareMatch  bool   `json:"hardware_match"`
	RemainingDays  int    `json:"remaining_days"`
	ActivationTime string `json:"activation_time"`
	IsValid        bool   `json:"is_valid"`
}
