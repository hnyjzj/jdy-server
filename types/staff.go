package types

import (
	"errors"
	"jdy/enums"

	"github.com/go-playground/validator/v10"
)

// 员工请求
type StaffReq struct {
	Platform PlatformType `json:"platform" binding:"required"` // 平台

	Account *StaffAccountReq `json:"account,omitempty"` // 账号信息
	WxWork  *StaffWxWorkReq  `json:"wxwork,omitempty"`  // 企业微信信息
}

func (req *StaffReq) Validate() error {
	validate := validator.New()
	switch req.Platform {
	case PlatformTypeAccount:
		if req.Account == nil {
			return errors.New("账号信息是必填项")
		}
		// 可以在这里对Account进行进一步的验证
		if err := validate.Struct(req.Account); err != nil {
			return err
		}
	case PlatformTypeWxWork:
		if req.WxWork == nil {
			return errors.New("企业微信信息是必填项")
		}

		// 可以在这里对WxWork进行进一步的验证
		if err := validate.Struct(req.WxWork); err != nil {
			return err
		}
	default:
		return errors.New("invalid Platform value")
	}
	return nil
}

// 员工账号请求
type StaffAccountReq struct {
	Phone    string `json:"phone" binding:"required,min=11,max=11,regex=^1\\d{10}$"` // 手机号
	Password string `json:"password" binding:"required"`                             // 密码

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

	Nickname string       `json:"nickname"`
	Avatar   string       `json:"avatar"`
	Email    string       `json:"email"`
	Gender   enums.Gender `json:"gender"`
}

// 更新请求
type StaffUpdateReq struct {
	Platform PlatformType `json:"platform" binding:"required"` // 平台

	Account *StaffUpdateAccountReq `json:"account,omitempty"` // 账号信息
	WxWork  *StaffUpdateWxWorkReq  `json:"wxwork,omitempty"`  // 企业微信信息
}

func (req *StaffUpdateReq) Validate() error {
	validate := validator.New()
	switch req.Platform {
	case PlatformTypeAccount:
		if req.Account == nil {
			return errors.New("账号信息是必填项")
		}
		// 可以在这里对Account进行进一步的验证
		if err := validate.Struct(req.Account); err != nil {
			return err
		}
	case PlatformTypeWxWork:
		if req.WxWork == nil {
			return errors.New("企业微信信息是必填项")
		}

		// 可以在这里对WxWork进行进一步的验证
		if err := validate.Struct(req.WxWork); err != nil {
			return err
		}
	default:
		return errors.New("invalid Platform value")
	}
	return nil
}

type StaffUpdateAccountReq struct {
	Password string `json:"password" ` // 密码

	Nickname string       `json:"nickname" binding:"min=2,max=50"` // 姓名
	Avatar   string       `json:"avatar"`                          // 头像
	Email    string       `json:"email" binding:"email"`           // 邮箱
	Gender   enums.Gender `json:"gender" binding:"oneof=0 1 2"`    // 性别
}

type StaffUpdateWxWorkReq struct {
	Code string `json:"code" binding:"required"`
}

type StaffWhere struct {
	Phone      string       `json:"phone" label:"手机号" show:"true" sort:"1" type:"string" input:"text"`
	Nickname   string       `json:"nickname" label:"姓名" show:"true" sort:"3" type:"string" input:"text"`
	Gender     enums.Gender `json:"gender" label:"性别" show:"true" sort:"4" type:"number" input:"select" preset:"typeMap"`
	IsDisabled bool         `json:"is_disabled" label:"是否禁用" show:"true" sort:"5" type:"boolean" input:"switch"`
}

type StaffListReq struct {
	PageReq
	Where StaffWhere `json:"where"`
}
