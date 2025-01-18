package enums

import (
	"errors"
)

/* 订单状态 */
// 全部、待付款、已取消、已完成、已退款、已收货、派维修、维修中、已维修、待取货、已预订、已核销
type OrderStatus int

const (
	OrderStatusAll          OrderStatus = iota // 全部
	OrderStatusWaitPay                         // 待付款
	OrderStatusCancel                          // 已取消
	OrderStatusComplete                        // 已完成
	OrderStatusRefund                          // 已退款
	OrderStatusReceived                        // 已收货
	OrderStatusSendRepair                      // 派维修
	OrderStatusRepairing                       // 维修中
	OrderStatusRepaired                        // 已维修
	OrderStatusWaitPickup                      // 待取货
	OrderStatusReserve                         // 已预订
	OrderStatusVerification                    // 已核销
)

var OrderStatusMap = map[OrderStatus]string{
	OrderStatusAll:          "全部",
	OrderStatusWaitPay:      "待付款",
	OrderStatusCancel:       "已取消",
	OrderStatusComplete:     "已完成",
	OrderStatusRefund:       "已退款",
	OrderStatusReceived:     "已收货",
	OrderStatusSendRepair:   "派维修",
	OrderStatusRepairing:    "维修中",
	OrderStatusRepaired:     "已维修",
	OrderStatusWaitPickup:   "待取货",
	OrderStatusReserve:      "已预订",
	OrderStatusVerification: "已核销",
}

var OrderStatusTypeMap = map[OrderType][]OrderStatus{
	OrderTypeSales:   {OrderStatusWaitPay, OrderStatusCancel, OrderStatusComplete, OrderStatusRefund},
	OrderTypeDeposit: {OrderStatusReserve, OrderStatusVerification, OrderStatusRefund},
	OrderTypeRepair:  {OrderStatusReceived, OrderStatusSendRepair, OrderStatusRepairing, OrderStatusWaitPickup, OrderStatusRepaired, OrderStatusComplete, OrderStatusCancel, OrderStatusRefund},
	OrderTypeOthers:  {OrderStatusWaitPay, OrderStatusCancel, OrderStatusComplete, OrderStatusRefund},
}

func (p OrderStatus) ToMap() any {
	return OrderStatusMap
}

func (p OrderStatus) InMap() error {
	if _, ok := OrderStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
