package config

import (
	"log"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type Configured struct {
	Server   Server   `mapstructure:"server"`
	JWT      JWT      `mapstructure:"jwt"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Storage  Storage  `mapstructure:"storage"`
	Wechat   Wechat   `mapstructure:"wechat"`
}

var Config *Configured

func Init() {
	// 设置默认值
	Config = new(Configured)
	defaults.SetDefaults(Config)
	// 读取配置文件
	viper.SetConfigName("env")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("没有读取到配置文件，将使用默认值。%v \n", err)
	}
	// 解析配置文件
	if err := viper.Unmarshal(&Config); err != nil {
		log.Printf("无法解码配置文件, %v \n", err)
	}
}
