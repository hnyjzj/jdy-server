package config

import (
	"fmt"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
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

type WechatService struct {
	JdyWork *work.Work
}

func NewWechatService() *WechatService {
	return &WechatService{
		JdyWork: newJdyWork(),
	}
}

// @see https://powerwechat.artisan-cloud.com/zh/wecom/
func newJdyWork() *work.Work {
	WeComApp, err := work.NewWork(&work.UserConfig{
		CorpID:  Config.Wechat.Work.CorpID,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: Config.Wechat.Work.Jdy.Id,     // 内部应用的app id
		Secret:  Config.Wechat.Work.Jdy.Secret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: Config.Wechat.Work.Jdy.Callback,
			Scopes:   []string{"snsapi_privateinfo"},
		},
		Log: work.Log{
			Level:  "debug",
			File:   "./logs/wechat/wxwork_info.log",
			Error:  "./logs/wechat/wxwork_error.log",
			Stdout: false, //  是否打印在终端
		},
		Cache: kernel.NewRedisClient(&kernel.UniversalOptions{
			Addrs:    []string{fmt.Sprintf("%s:%d", Config.Redis.Host, Config.Redis.Port)},
			Password: Config.Redis.Password,
			DB:       Config.Redis.Db + 1,
		}),
		HttpDebug: false,
	})

	if err != nil {
		panic(err)
	}

	return WeComApp
}
