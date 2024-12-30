package types

import (
	"jdy/enums"
)

type MemberWhere struct {
	Phone       *string      `json:"phone" label:"手机号" show:"true" sort:"1" type:"string" input:"text" required:"true"`
	Name        string       `json:"name" label:"姓名" show:"true" sort:"2" type:"string" input:"text"`
	Gender      enums.Gender `json:"gender" label:"性别" show:"true" sort:"3" type:"number" input:"select" preset:"typeMap"`
	Birthday    string       `json:"birthday" label:"生日" show:"true" sort:"4" type:"date" input:"date"`
	Anniversary string       `json:"anniversary" label:"纪念日" show:"true" sort:"5" type:"date" input:"date"`
	Nickname    string       `json:"nickname" label:"昵称" show:"true" sort:"6" type:"string" input:"text"`

	Level      enums.MemberLevel `json:"level" label:"等级" show:"true" sort:"8" type:"number" input:"select" preset:"typeMap"`
	Integral   float64           `json:"integral" label:"积分" show:"true" sort:"9" type:"number" input:"text"`
	BuyCount   int               `json:"buy_count" label:"购买次数" show:"true" sort:"10" type:"number" input:"text"`
	EventCount int               `json:"event_count" label:"活动次数" show:"true" sort:"11" type:"number" input:"text"`

	Source       enums.MemberSource `json:"source" label:"来源" show:"true" sort:"12" type:"number" input:"select" preset:"typeMap"`
	ConsultantId string             `json:"consultant_id" label:"顾问" show:"true" sort:"13" type:"string" input:"text"`
	StoreId      string             `json:"store_id" label:"门店" show:"true" sort:"14" type:"string" input:"text"`

	Status enums.MemberStatus `json:"status" label:"状态" show:"true" sort:"15" type:"number" input:"select" preset:"typeMap"`
}

type MemberCreateReq struct {
	Phone       *string      `json:"phone" binding:"required,min=11,max=11,regex=^1\\d{10}$"`
	Name        string       `json:"name" binding:"required"`
	Gender      enums.Gender `json:"gender" binding:"oneof=0 1 2"`
	Birthday    string       `json:"birthday" binding:"-"`
	Anniversary string       `json:"anniversary" binding:"-"`
	Nickname    string       `json:"nickname" binding:"-"`
	IDCard      string       `json:"id_card" binding:"-"`

	ConsultantId string `json:"consultant_id" binding:"-"`
	StoreId      string `json:"store_id" binding:"-"`
}

type MemberUpdateReq struct {
	Id string `json:"id" binding:"required"`
	MemberCreateReq
}

type MemberListReq struct {
	PageReq
	Where MemberWhere `json:"where" binding:"required"`
}

type MemberInfoReq struct {
	Id string `json:"id" binding:"required"`
}
