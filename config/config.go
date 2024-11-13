package config

type Config struct {
	System    System    `mapstructure:"system"`
	LogConfig LogConfig `mapstructure:"log"`
}
