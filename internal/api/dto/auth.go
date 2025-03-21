package dto

type Login struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	IgnoreCaptcha bool   `json:"ignoreCaptcha"`
	Captcha       string `json:"captcha"`
	CaptchaID     string `json:"captchaID"`
	AuthMethod    string `json:"authMethod"`
}

type CaptchaResponse struct {
	CaptchaID string `json:"captchaID"`
	ImagePath string `json:"imagePath"`
}

type UserLoginInfo struct {
	Name      string `json:"name"`
	Token     string `json:"token"`
	MfaStatus string `json:"mfaStatus"`
}
