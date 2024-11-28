package configs

type System struct {
	Port        string `mapstructure:"port"`
	BaseDir     string `mapstructure:"base_dir"`
	Mode        string `mapstructure:"mode"`
	Version     string `mapstructure:"version"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Entrance    string `mapstructure:"entrance"`
	LogPath     string `mapstructure:"log_path"`
	DbPath      string `mapstructure:"db_path"`
	DbFile      string `mapstructure:"db_file"`
	DataDir     string `mapstructure:"data_dir"`
	BindAddress string `mapstructure:"bindAddress"`
	EncryptKey  string `mapstructure:"encrypt_key"`
}
