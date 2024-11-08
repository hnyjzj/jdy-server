package commonlogic

import (
	"github.com/mojocn/base64Captcha"
)

type CaptchaLogic struct{}

type CaptchaRes struct {
	Id   string `json:"id"`
	Code string `json:"code"`
}

// 获取图片验证码
func (l *CaptchaLogic) ImageCaptcha() (*CaptchaRes, error) {
	var (
		res = &CaptchaRes{}
		err error
	)
	// 生成图片验证码
	res.Id, res.Code, err = l.CreateImage()
	if err != nil {
		return nil, err
	}

	return res, nil
}

var result = base64Captcha.DefaultMemStore

// 生成图片验证码
func (l *CaptchaLogic) CreateImage() (string, string, error) {

	driver := &base64Captcha.DriverDigit{
		Height:   50,
		Width:    100,
		Length:   5,
		MaxSkew:  0.15,
		DotCount: 10,
	}

	// 创建验证码并传入创建的类型的配置，以及存储的对象
	c := base64Captcha.NewCaptcha(driver, result)
	id, b64s, _, err := c.Generate()

	return id, b64s, err
}

// 对比验证码
func (l *CaptchaLogic) VerifyCaptcha(id string, VerifyValue string) bool {
	var result = base64Captcha.DefaultMemStore

	return result.Verify(id, VerifyValue, true)
}
