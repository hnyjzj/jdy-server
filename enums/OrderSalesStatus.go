package enums

import (
	"errors"
)

/* 销售单状态 */
// 待付款、已取消、已完成、有退货、已退
type OrderSalesStatus int

const (
	OrderSalesStatusWaitPay  OrderSalesStatus = iota + 1 // 待付款
	OrderSalesStatusCancel                               // 已取消
	OrderSalesStatusComplete                             // 已完成
	OrderSalesStatusRefund                               // 有退货
	OrderSalesStatusReturn                               // 已退
)

var OrderSalesStatusMap = map[OrderSalesStatus]string{
	OrderSalesStatusWaitPay:  "待付款",
	OrderSalesStatusCancel:   "已取消",
	OrderSalesStatusComplete: "已完成",
	OrderSalesStatusRefund:   "有退货",
	OrderSalesStatusReturn:   "已退",
}

func (p OrderSalesStatus) ToMap() any {
	return OrderSalesStatusMap
}

func (p OrderSalesStatus) InMap() error {
	if _, ok := OrderSalesStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
