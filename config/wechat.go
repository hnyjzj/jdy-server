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

	Jdy      Agent `mapstructure:"jdy"`      // 应用
	Contacts Agent `mapstructure:"contacts"` // 通讯录
}

type Agent struct {
	Id     int    `mapstructure:"id"`     // 应用ID
	Secret string `mapstructure:"secret"` // 应用秘钥
	Home   string `mapstructure:"home"`   // 应用首页

	Token  string `mapstructure:"token"`  // 应用Token
	AESKey string `mapstructure:"aeskey"` // 应用AesKey

	CallbackURL   string `mapstructure:"callback_url" default:"https://example.cn/callback/wxwork"`         // 回调地址
	CallbackOAuth string `mapstructure:"callback_oauth" default:"https://example.cn/callback/wxwork/oauth"` // 回调地址
}

type WechatService struct {
	JdyWork      *work.Work
	ContactsWork *work.Work
}

func NewWechatService() *WechatService {
	return &WechatService{
		JdyWork:      newJdyWork(),
		ContactsWork: newContactsWork(),
	}
}

// @see https://powerwechat.artisan-cloud.com/zh/wecom/
func newJdyWork() *work.Work {
	app := Config.Wechat.Work
	WeComApp, err := work.NewWork(&work.UserConfig{
		CorpID:      app.CorpID,     // 企业微信的app id，所有企业微信共用一个。
		AgentID:     app.Jdy.Id,     // 内部应用的app id
		Secret:      app.Jdy.Secret, // 内部应用的app secret
		Token:       app.Jdy.Token,
		AESKey:      app.Jdy.AESKey,
		CallbackURL: app.Jdy.CallbackURL,
		OAuth: work.OAuth{
			Callback: app.Jdy.CallbackOAuth,
			Scopes:   []string{"snsapi_privateinfo"},
		},
		Log: work.Log{
			Level:  "debug",
			File:   "./logs/wechat/wxwork_info_jdy.log",
			Error:  "./logs/wechat/wxwork_error_jdy.log",
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
func newContactsWork() *work.Work {
	app := Config.Wechat.Work
	WeComApp, err := work.NewWork(&work.UserConfig{
		CorpID: app.CorpID,          // 企业微信的app id，所有企业微信共用一个。
		Secret: app.Contacts.Secret, // 通讯录的secret
		OAuth: work.OAuth{
			Callback: app.Contacts.CallbackOAuth,
			Scopes:   []string{"snsapi_privateinfo"},
		},
		Log: work.Log{
			Level:  "debug",
			File:   "./logs/wechat/wxwork_info_contacts.log",
			Error:  "./logs/wechat/wxwork_error_contacts.log",
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
