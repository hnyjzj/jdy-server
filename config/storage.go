package config

const (
	StorageTypeLocal = "local"
)

type Storage struct {
	// 默认方式
	Type string `mapstructure:"default" default:"local"`
	// 文件前缀
	Prefix string `mapstructure:"prefix" default:"uploads"`

	// 本地存储
	Local Local `mapstructure:"local"`
}

type Local struct {
	// 本地存储路径
	Root string `mapstructure:"root" default:"./uploads/"`
}
