package enums

import "errors"

/* 黄金价格状态 */
// 待审批、已审批、已驳回
type GoldPriceStatus int

const (
	GoldPriceStatusPending  GoldPriceStatus = iota // 待审批
	GoldPriceStatusApproved                        // 已审批
	GoldPriceStatusRejected                        // 已驳回
)

var GoldPriceStatusMap = map[GoldPriceStatus]string{
	GoldPriceStatusPending:  "待审批",
	GoldPriceStatusApproved: "已审批",
	GoldPriceStatusRejected: "已驳回",
}

func (p GoldPriceStatus) ToMap() any {
	return GoldPriceStatusMap
}

func (p GoldPriceStatus) InMap() error {
	if _, ok := GoldPriceStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
