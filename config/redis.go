package config

type Redis struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host" default:"127.0.0.1"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port" default:"6379"`
	Password string `mapstructure:"password" json:"password" yaml:"password" default:""`
	Db       int    `mapstructure:"db" json:"db" yaml:"db" default:"0"`
}
