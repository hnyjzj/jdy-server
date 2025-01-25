package enums

import "errors"

/* 会员积分变更类型 */
// 全部、消费、充值、兑换、退款、取消兑换、取消退款、人工调整
type MemberIntegralChangeType int

const (
	MemberIntegralChangeTypeAll            MemberIntegralChangeType = iota // 全部
	MemberIntegralChangeTypeConsume                                        // 消费
	MemberIntegralChangeTypeRecharge                                       // 充值
	MemberIntegralChangeTypeExchange                                       // 兑换
	MemberIntegralChangeTypeRefund                                         // 退款
	MemberIntegralChangeTypeCancelExchange                                 // 取消兑换
	MemberIntegralChangeTypeCancelRefund                                   // 取消退款
	MemberIntegralChangeTypeAdjust                                         // 人工调整
)

var MemberIntegralChangeTypeMap = map[MemberIntegralChangeType]string{
	MemberIntegralChangeTypeAll:            "全部",
	MemberIntegralChangeTypeConsume:        "消费",
	MemberIntegralChangeTypeRecharge:       "充值",
	MemberIntegralChangeTypeExchange:       "兑换",
	MemberIntegralChangeTypeRefund:         "退款",
	MemberIntegralChangeTypeCancelExchange: "取消兑换",
	MemberIntegralChangeTypeCancelRefund:   "取消退款",
	MemberIntegralChangeTypeAdjust:         "人工调整",
}

func (p MemberIntegralChangeType) ToMap() any {
	return MemberIntegralChangeTypeMap
}

func (p MemberIntegralChangeType) InMap() error {
	if _, ok := MemberIntegralChangeTypeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
