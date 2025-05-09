package enums

import (
	"errors"
)

/* 其他单状态 */
// 待付款、已取消、已完成、有退款
type OrderOtherStatus int

const (
	OrderOtherStatusWaitPay  OrderOtherStatus = iota + 1 // 待付款
	OrderOtherStatusCancel                               // 已取消
	OrderOtherStatusComplete                             // 已完成
	OrderOtherStatusRefund                               // 有退款
)

var OrderOtherStatusMap = map[OrderOtherStatus]string{
	OrderOtherStatusWaitPay:  "待付款",
	OrderOtherStatusCancel:   "已取消",
	OrderOtherStatusComplete: "已完成",
	OrderOtherStatusRefund:   "有退款",
}

func (p OrderOtherStatus) ToMap() any {
	return OrderOtherStatusMap
}

func (p OrderOtherStatus) InMap() error {
	if _, ok := OrderOtherStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
