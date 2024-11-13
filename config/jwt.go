package config

type JWT struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret" default:"secret"` // 密钥

	Expire int64 `mapstructure:"expire" json:"expire" yaml:"expire" default:"7200"` // 过期时间(秒)
}
