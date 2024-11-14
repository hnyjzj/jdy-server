package errors

var (
	// 参数错误
	ErrInvalidParam = New("参数错误", C400)
	// 验证码错误
	ErrInvalidCaptcha = New("验证码错误", C403)
)
