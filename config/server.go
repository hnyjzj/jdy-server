package config

type Server struct {
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode" default:"release"` // 服务运行的模式 debug、test、release
	Port int    `mapstructure:"port" json:"port" yaml:"port" default:"9527" `
}
