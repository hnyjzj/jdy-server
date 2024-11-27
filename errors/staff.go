package errors

var (
	// 账号不存在
	ErrStaffNotFound = New("账号不存在", C404)
	// 账号已被禁用
	ErrStaffDisabled = New("账号已被禁用", C403)
	// 密码错误
	ErrPasswordIncorrect = New("密码错误", C403)
)
