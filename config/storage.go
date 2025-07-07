package config

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

const (
	StorageTypeLocal        = "local"
	StorageTypeTencentCloud = "tencent_cloud"
)

type Storage struct {
	// 默认方式
	Type string `mapstructure:"type" default:"local"`
	// 文件前缀
	Prefix string `mapstructure:"prefix" default:"uploads"`

	// 本地存储
	Local Local `mapstructure:"local"`
	// 腾讯云存储
	TencentCloud TencentCloud `mapstructure:"tencent_cloud"`
}

// 本地存储
type Local struct {
	Root string `mapstructure:"root" default:"./uploads/"` // 根目录
}

// 腾讯云存储
type TencentCloud struct {
	Root string `mapstructure:"root" default:"uploads"` // 根目录

	Region    string `mapstructure:"region" default:"ap-guangzhou"` // 地域
	Bucket    string `mapstructure:"bucket" default:"jdy-oss"`      // 存储桶
	SecretId  string `mapstructure:"secret_id"`                     // 密钥ID
	SecretKey string `mapstructure:"secret_key"`                    // 密钥Key

	Timeout time.Duration `mapstructure:"timeout" default:"100s"` // 超时时间
}

func NewTencentCloudStorage() *cos.Client {
	// 初始化cos客户端
	urlStr, _ := url.Parse(fmt.Sprintf(
		"https://%s.cos.%s.myqcloud.com",
		Config.Storage.TencentCloud.Bucket,
		Config.Storage.TencentCloud.Region,
	))
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Timeout: Config.Storage.TencentCloud.Timeout, // 超时时间
		Transport: &cos.AuthorizationTransport{
			SecretID:  Config.Storage.TencentCloud.SecretId,
			SecretKey: Config.Storage.TencentCloud.SecretKey,
		},
	})
	return client
}
