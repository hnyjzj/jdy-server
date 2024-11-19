package errors

var (
	// 参数错误
	ErrInvalidParam = New("参数错误", C400)
	// 结果不存在
	ErrNotFound = New("结果不存在", C404)
)

var (
	// 验证码错误
	ErrInvalidCaptcha = New("验证码错误", C403)
)
