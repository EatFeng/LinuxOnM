package dto

type SettingInfo struct {
	UserName       string `json:"userName"`
	Email          string `json:"email"`
	SystemIP       string `json:"systemIP"`
	SystemVersion  string `json:"systemVersion"`
	DockerSockPath string `json:"dockerSockPath"`
	DeveloperMode  string `json:"developerMode"`

	SessionTimeout string `json:"sessionTimeout"`
	LocalTime      string `json:"localTime"`
	TimeZone       string `json:"timeZone"`
	NtpSite        string `json:"ntpSite"`

	Port           string `json:"port"`
	Ipv6           string `json:"ipv6"`
	BindAddress    string `json:"bindAddress"`
	Theme          string `json:"theme"`
	MenuTabs       string `json:"menuTabs"`
	DefaultNetwork string `json:"defaultNetwork"`
	LastCleanTime  string `json:"lastCleanTime"`
	LastCleanSize  string `json:"lastCleanSize"`
	LastCleanData  string `json:"lastCleanData"`

	ServerPort             string `json:"serverPort"`
	SSL                    string `json:"ssl"`
	SSLType                string `json:"sslType"`
	BindDomain             string `json:"bindDomain"`
	AllowIPs               string `json:"allowIPs"`
	SecurityEntrance       string `json:"securityEntrance"`
	ExpirationDays         string `json:"expirationDays"`
	ExpirationTime         string `json:"expirationTime"`
	ComplexityVerification string `json:"complexityVerification"`

	MonitorStatus    string `json:"monitorStatus"`
	MonitorInterval  string `json:"monitorInterval"`
	MonitorStoreDays string `json:"monitorStoreDays"`

	MessageType string `json:"messageType"`
	EmailVars   string `json:"emailVars"`

	FileRecycleBin string `json:"fileRecycleBin"`

	SnapshotIgnore string `json:"snapshotIgnore"`

	ProxyUrl        string `json:"proxyUrl"`
	ProxyType       string `json:"proxyType"`
	ProxyPort       string `json:"proxyPort"`
	ProxyUser       string `json:"proxyUser"`
	ProxyPasswd     string `json:"proxyPasswd"`
	ProxyPasswdKeep string `json:"proxyPasswdKeep"`
}

type SettingUpdate struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value"`
}
