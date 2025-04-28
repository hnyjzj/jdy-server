package enums

import (
	"errors"
)

/* 订金单状态 */
// 待付款、已取消、预定中、已核销、有退款、已退
type OrderDepositStatus int

const (
	OrderDepositStatusWaitPay  OrderDepositStatus = iota + 1 // 待付款
	OrderDepositStatusCancel                                 // 已取消
	OrderDepositStatusBooking                                // 预定中
	OrderDepositStatusComplete                               // 已核销
	OrderDepositStatusRefund                                 // 有退款
	OrderDepositStatusReturn                                 // 已退
)

var OrderDepositStatusMap = map[OrderDepositStatus]string{
	OrderDepositStatusWaitPay:  "待付款",
	OrderDepositStatusCancel:   "已取消",
	OrderDepositStatusBooking:  "预定中",
	OrderDepositStatusComplete: "已核销",
	OrderDepositStatusRefund:   "有退款",
	OrderDepositStatusReturn:   "已退",
}

func (p OrderDepositStatus) ToMap() any {
	return OrderDepositStatusMap
}

func (p OrderDepositStatus) InMap() error {
	if _, ok := OrderDepositStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
