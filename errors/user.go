package errors

var (
	// 用户不存在
	ErrStaffNotFound = New("用户不存在", C404)
	// 用户已被禁用
	ErrStaffDisabled = New("用户已被禁用", C403)
	// 密码错误
	ErrPasswordIncorrect = New("密码错误", C403)
)
