package authtype

type LoginReq struct {
	Phone    string `json:"phone" binding:"required,regex=^1\\d{10}$"`
	Password string `json:"password" binding:"required"`

	CaptchaId string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}

type LoginOAuthReq struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

type LoginOAuthRes struct {
	Token string `json:"token"`
	Res   any    `json:"res"`
}
