package dto

type Login struct {
	Name          string `json:"name"`
	Password      string `json:"password"`
	IgnoreCapycha bool   `json:"ignore_capycha"`
	Captcha       string `json:"captcha"`
	CaptchaID     string `json:"captchaID"`
	AuthMethod    string `json:"authMethod"`
}
