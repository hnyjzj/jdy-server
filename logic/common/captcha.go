package commonlogic

import (
	"context"
	"jdy/config"
	"jdy/service/redis"
	"time"

	"github.com/mojocn/base64Captcha"
)

type CaptchaLogic struct{}

type CaptchaRes struct {
	Id   string `json:"id"`
	Code string `json:"code"`
	Val  string `json:"val,omitempty"`
}

// 获取图片验证码
func (l *CaptchaLogic) ImageCaptcha() (*CaptchaRes, error) {
	var (
		res    = &CaptchaRes{}
		err    error
		answer string
	)
	// 生成图片验证码
	res.Id, res.Code, answer, err = l.CreateImage()
	if err != nil {
		return nil, err
	}

	if config.Config.Server.Mode == "debug" {
		res.Val = answer
	}

	return res, nil
}

var stroe = MyStore{
	ctx: context.Background(),
}

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
	c := base64Captcha.NewCaptcha(driver, &stroe)
	id, b64s, answer, err := c.Generate()

	return id, b64s, answer, err
}

// 对比验证码
func (l *CaptchaLogic) VerifyCaptcha(id string, VerifyValue string) bool {
	return stroe.Verify(id, VerifyValue, true)
}

// 验证码存储
type MyStore struct {
	base64Captcha.Store
	ctx context.Context
}

var prefix = "captcha_"

// 设置验证码
func (s *MyStore) Set(id string, value string) error {
	// 存到redis，设置过期时间为n秒
	result := redis.Client.Set(s.ctx, prefix+id, value, time.Second*time.Duration(120))
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// 获取验证码
func (s *MyStore) Get(id string, clear bool) string {
	result := redis.Client.Get(s.ctx, prefix+id)
	if result.Err() != nil {
		return ""
	}

	if clear {
		redis.Client.Del(s.ctx, prefix+id)
	}

	return result.Val()
}

// 验证验证码
func (MyStore) Verify(id, value string, clear bool) bool {
	result := redis.Client.Get(stroe.ctx, prefix+id)
	if result.Err() != nil {
		return false
	}

	if clear {
		redis.Client.Del(stroe.ctx, prefix+id)
	}

	if result.Val() != value {
		return false
	}

	return true
}
