package types

import (
	"errors"
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type MemberWhere struct {
	Phone       string       `json:"phone" label:"手机号" find:"true" sort:"1" type:"string" input:"text"`
	Name        string       `json:"name" label:"姓名" find:"true" sort:"2" type:"string" input:"text"`
	Gender      enums.Gender `json:"gender" label:"性别" find:"true" sort:"3" type:"number" input:"select" preset:"typeMap"`
	Birthday    string       `json:"birthday" label:"生日" find:"true" sort:"4" type:"date" input:"date"`
	Anniversary string       `json:"anniversary" label:"纪念日" find:"true" sort:"5" type:"date" input:"date"`
	Nickname    string       `json:"nickname" label:"昵称" find:"true" sort:"6" type:"string" input:"text"`

	Level      enums.MemberLevel `json:"level" label:"等级" find:"true" sort:"8" type:"number" input:"select" preset:"typeMap"`
	Integral   decimal.Decimal   `json:"integral" label:"积分" find:"true" sort:"9" type:"number" input:"text"`
	BuyCount   int               `json:"buy_count" label:"购买次数" find:"true" sort:"10" type:"number" input:"text"`
	EventCount int               `json:"event_count" label:"活动次数" find:"true" sort:"11" type:"number" input:"text"`

	Source       enums.MemberSource `json:"source" label:"来源" find:"true" sort:"12" type:"number" input:"select" preset:"typeMap"`
	ConsultantId string             `json:"consultant_id" label:"顾问" find:"true" sort:"13" type:"string" input:"text"`
	StoreId      string             `json:"store_id" label:"门店" find:"true" sort:"14" type:"string" input:"text" required:"true"`

	Status         enums.MemberStatus `json:"status" label:"状态" find:"true" sort:"15" type:"number" input:"select" preset:"typeMap"`
	ExternalUserId string             `json:"external_user_id" label:"外部用户id" find:"true" sort:"16" type:"string" input:"text"`
}

type MemberCreateReq struct {
	Phone       *string      `json:"phone" binding:"required,min=11,max=11,regex=^1\\d{10}$"`
	Name        string       `json:"name"`
	Gender      enums.Gender `json:"gender"`
	Birthday    string       `json:"birthday"`
	Anniversary string       `json:"anniversary"`
	Nickname    string       `json:"nickname"`
	IDCard      string       `json:"id_card"`

	ConsultantId   string `json:"consultant_id" binding:"required"`
	StoreId        string `json:"store_id" binding:"required"`
	ExternalUserId string `json:"external_user_id"`
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
	Id             string `json:"id" binding:"-"`
	ExternalUserId string `json:"external_user_id" binding:"-"`
}

type MemberConsumptionsReq struct {
	Id string `json:"id" binding:"required"`
}

func (req *MemberInfoReq) Validate() error {
	if req.Id == "" && req.ExternalUserId == "" {
		return errors.New("用户标识不能同时为空")
	}

	return nil
}
