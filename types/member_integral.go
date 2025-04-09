package types

import (
	"jdy/enums"

	"github.com/shopspring/decimal"
)

type MemberIntegralListReq struct {
	PageReq
	Where MemberIntegralWhere `json:"where" binding:"required"`
}

type MemberIntegralWhere struct {
	MemberId   string                         `json:"member_id" label:"会员id" find:"true" sort:"1" type:"string" input:"text" binding:"required" `
	ChangeType enums.MemberIntegralChangeType `json:"change_type" label:"变更类型" find:"true" sort:"2" type:"number" input:"select" preset:"typeMap"`
}

type MemberIntegralChangeReq struct {
	MemberId string          `json:"id" binding:"required"`
	Change   decimal.Decimal `json:"change" binding:"required"`
	Remark   string          `json:"remark" binding:"required"`
}

type MemberIntegralRuleReq struct {
	Class int `json:"class" binding:"required"`
}
