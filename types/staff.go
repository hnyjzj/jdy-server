package types

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// 员工请求
type StaffReq struct {
	Platform PlatformType `json:"platform" binding:"required"` // 平台

	Account *StaffAccountReq `json:"account,omitempty"` // 账号信息
	WxWork  *StaffWxWorkReq  `json:"wxwork,omitempty"`  // 企业微信信息
}

func (StaffReq) ValidateStaffReq(staffReq *StaffReq) error {
	validate := validator.New()
	switch staffReq.Platform {
	case PlatformTypeAccount:
		if staffReq.Account == nil {
			return errors.New("账号信息是必填项")
		}
		// 可以在这里对Account进行进一步的验证
		if err := validate.Struct(staffReq.Account); err != nil {
			return err
		}
	case PlatformTypeWxWork:
		if staffReq.WxWork == nil {
			return errors.New("企业微信信息是必填项")
		}

		// 可以在这里对WxWork进行进一步的验证
		if err := validate.Struct(staffReq.WxWork); err != nil {
			return err
		}
	default:
		return errors.New("invalid Platform value")
	}
	return nil
}

// 员工账号请求
type StaffAccountReq struct {
	Username string `json:"username" binding:"required,min=2,max=50,regex=^[a-zA-Z0-9_]+$"` // 用户名
	Phone    string `json:"phone" binding:"required,min=11,max=11,regex=^1\\d{10}$"`        // 手机号
	Password string `json:"password" binding:"required"`                                    // 密码

	Nickname string `json:"nickname" binding:"required,min=2,max=50,regex=^[\u4e00-\u9fa5]+$"` // 姓名
	Avatar   string `json:"avatar"`                                                            // 头像
	Email    string `json:"email"`                                                             // 邮箱
}

// 企业微信信息
type StaffWxWorkReq struct {
	UserId []string `json:"user_id" binding:"required"` // 用户ID
}

// 员工响应
type StaffRes struct {
	Phone string `json:"phone"`

	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Gender   uint   `json:"gender"`
}
