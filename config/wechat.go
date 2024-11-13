package config

import (
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
)

type Wechat struct {
	Work WechatWork `mapstructure:"work"` // 企业微信
}

type WechatWork struct {
	CorpID string `mapstructure:"corp_id"` // 企业ID
	Secret string `mapstructure:"secret"`  // 通讯录秘钥

	Jdy agent `mapstructure:"jdy"` // 应用
}

type agent struct {
	Id     int    `mapstructure:"id"`     // 应用ID
	Secret string `mapstructure:"secret"` // 应用秘钥

	Token  string `mapstructure:"Token"`  // 应用Token
	AesKey string `mapstructure:"AesKey"` // 应用AesKey

	Callback string `mapstructure:"callback" default:"https://example.cn/wxwork/callback"` // 回调地址
}

var JdyAgent *work.Work

func WechatHandler() {
	JdyAgent = NewJdyAgent()
}

func NewJdyAgent() *work.Work {
	WeComApp, err := work.NewWork(&work.UserConfig{
		CorpID:  Config.Wechat.Work.CorpID,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: Config.Wechat.Work.Jdy.Id,     // 内部应用的app id
		Secret:  Config.Wechat.Work.Jdy.Secret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: Config.Wechat.Work.Jdy.Callback,
			Scopes:   []string{"snsapi_privateinfo"},
		},
		Log: work.Log{
			Level:  "info",
			File:   "./temp/logs/wechat.log",
			Stdout: false, //  是否打印在终端
		},
		HttpDebug: false,
	})

	if err != nil {
		panic(err)
	}

	return WeComApp
}
