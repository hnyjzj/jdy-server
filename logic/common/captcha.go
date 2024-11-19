package common

import (
	"context"
	"jdy/logic"
	"jdy/service/redis"
	"time"

	"github.com/mojocn/base64Captcha"
)

type CaptchaLogic struct {
	logic.Base
}

type CaptchaRes struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Val  string `json:"val,omitempty"`
}

// 获取图片验证码
func (l *CaptchaLogic) ImageCaptcha() (*CaptchaRes, error) {
	var (
		res = &CaptchaRes{}
		err error
	)

	// 生成图片验证码
	res.Id, res.Code, _, err = l.CreateImage()
	if err != nil {
		return nil, err
	}

	return res, nil
}

var (
	store = &MyStore{}
	ctx   = context.Background()
)

// 生成图片验证码
func (l *CaptchaLogic) CreateImage() (string, string, string, error) {
	driver := &base64Captcha.DriverDigit{
		Height:   50,
		Width:    100,
		Length:   5,
		MaxSkew:  0.15,
		DotCount: 10,
	}

	// 创建验证码并传入创建的类型的配置，以及存储的对象
	c := base64Captcha.NewCaptcha(driver, store)
	id, b64s, answer, err := c.Generate()

	return id, b64s, answer, err
}

// 对比验证码
func (l *CaptchaLogic) VerifyCaptcha(id string, VerifyValue string) bool {
	return store.Verify(id, VerifyValue, true)
}

// 验证码存储
type MyStore struct {
	base64Captcha.Store
}

// 设置验证码
func (s *MyStore) Set(id string, value string) error {
	key := GetRedisKey(id)
	// 存到redis，设置过期时间为n秒
	result := redis.Client.Set(ctx, key, value, time.Second*time.Duration(120))
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// 获取验证码
func (s MyStore) Get(id string, clear bool) string {
	key := GetRedisKey(id)
	result := redis.Client.Get(ctx, key)
	if result.Err() != nil {
		return ""
	}

	if clear {
		redis.Client.Del(ctx, key)
	}

	return result.Val()
}

// 验证验证码
func (MyStore) Verify(id, value string, clear bool) bool {
	key := GetRedisKey(id)
	result := redis.Client.Get(ctx, key)
	if result.Err() != nil {
		return false
	}

	if clear {
		redis.Client.Del(ctx, key)
	}

	if result.Val() != value {
		return false
	}

	return true
}

// 获取 redis 名称
func GetRedisKey(id string) string {
	return "captcha_" + id
}
