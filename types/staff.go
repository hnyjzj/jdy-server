package types

import (
	"jdy/enums"
)

// 员工请求
type StaffReq struct {
	Username string `json:"username" binding:"required"`                             // 用户名
	Phone    string `json:"phone" binding:"required,min=11,max=11,regex=^1\\d{10}$"` // 手机号
	Password string `json:"password" binding:"required"`                             // 密码

	Nickname string       `json:"nickname" binding:"required,min=2,max=50,regex=^[\u4e00-\u9fa5]+$"` // 姓名
	Avatar   string       `json:"avatar"`                                                            // 头像
	Email    string       `json:"email"`                                                             // 邮箱
	Gender   enums.Gender `json:"gender"`                                                            // 性别
}

// 员工响应
type StaffRes struct {
	Id    string `json:"id"`
	Phone string `json:"phone"`

	Nickname string       `json:"nickname"`
	Avatar   string       `json:"avatar"`
	Email    string       `json:"email"`
	Gender   enums.Gender `json:"gender"`
}

// 更新请求
type StaffUpdateReq struct {
	Code     string `json:"code"`      // 授权码
	Password string `json:"password" ` // 密码

	Nickname string       `json:"nickname" binding:"min=2,max=50"` // 姓名
	Avatar   string       `json:"avatar"`                          // 头像
	Email    string       `json:"email" binding:"email"`           // 邮箱
	Gender   enums.Gender `json:"gender" binding:"oneof=0 1 2"`    // 性别
}

type StaffWhere struct {
	Username   string       `json:"username" label:"用户名" find:"true" create:"true" required:"true" sort:"1" type:"string" input:"text"`
	Phone      string       `json:"phone" label:"手机号" find:"true" create:"true" required:"true" sort:"2" type:"string" input:"text"`
	Nickname   string       `json:"nickname" label:"姓名" find:"true" create:"true" required:"true" sort:"3" type:"string" input:"text"`
	Email      string       `json:"email" label:"邮箱" find:"true"  create:"true" sort:"3" type:"string" input:"text"`
	Gender     enums.Gender `json:"gender" label:"性别" find:"true" create:"true" sort:"4" type:"number" input:"select" preset:"typeMap"`
	Avatar     string       `json:"avatar" label:"头像" create:"true" sort:"4" type:"string" input:"upload"`
	Password   string       `json:"password" label:"密码" create:"true" sort:"4" type:"string" input:"password"`
	IsDisabled bool         `json:"is_disabled" label:"是否禁用" find:"true" create:"true" sort:"5" type:"boolean" input:"switch"`
}

type StaffListReq struct {
	PageReq
	Where StaffWhere `json:"where"`
}

type StaffInfoReq struct {
	Id string `json:"id" binding:"required"`
}
