package logic_error

import "net/http"

var (
	// 用户不存在
	ErrUserNotFound = New(http.StatusNotFound, "用户不存在")
	// 用户已被禁用
	ErrUserDisabled = New(http.StatusForbidden, "用户已被禁用")
	// 密码错误
	ErrPasswordIncorrect = New(http.StatusBadRequest, "密码错误")
)
