package errors

var (
	// 员工不存在
	ErrStaffNotFound = New("员工不存在", C404)
	// 无权访问
	ErrStaffUnauthorized = New("无权访问", C403)
	// 员工已被禁用
	ErrStaffDisabled = New("员工已被禁用", C403)
	// 密码错误
	ErrPasswordIncorrect = New("密码错误", C403)
)
