package types

type LoginReq struct {
	Phone    string `json:"phone" binding:"required,regex=^1\\d{10}$"` // 手机号
	Password string `json:"password" binding:"required"`               // 密码

	CaptchaId string `json:"captcha_id" binding:"required"` // 验证码ID
	Captcha   string `json:"captcha" binding:"required"`    // 验证码
}

type LoginOAuthReq struct {
	Code  string `json:"code" binding:"required"`  // 授权码
	State string `json:"state" binding:"required"` // 状态码
}
