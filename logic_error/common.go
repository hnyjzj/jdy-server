package logic_error

import (
	"net/http"
)

var (
	// 参数错误
	ErrInvalidParam = New(http.StatusBadRequest, "参数错误")
	// 验证码错误
	ErrInvalidCaptcha = New(http.StatusForbidden, "验证码错误")
)
