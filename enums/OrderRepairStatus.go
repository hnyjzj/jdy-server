package enums

import (
	"errors"
)

/* 维修单状态 */
// 待付款、已取消、门店已收货、已送出维修、维修中、已维修送回、待取货、已完成
type OrderRepairStatus int

const (
	OrderRepairStatusWaitPay       OrderRepairStatus = iota + 1 // 待付款
	OrderRepairStatusCancel                                     // 已取消
	OrderRepairStatusStoreReceived                              // 门店已收货
	OrderRepairStatusSendOut                                    // 已送出维修
	OrderRepairStatusRepairing                                  // 维修中
	OrderRepairStatusSendBack                                   // 已维修送回
	OrderRepairStatusWaitPickUp                                 // 待取货
	OrderRepairStatusComplete                                   // 已完成
)

var OrderRepairStatusMap = map[OrderRepairStatus]string{
	OrderRepairStatusWaitPay:       "待付款",
	OrderRepairStatusCancel:        "已取消",
	OrderRepairStatusStoreReceived: "门店已收货",
	OrderRepairStatusSendOut:       "已送出维修",
	OrderRepairStatusRepairing:     "维修中",
	OrderRepairStatusSendBack:      "已维修送回",
	OrderRepairStatusWaitPickUp:    "待取货",
	OrderRepairStatusComplete:      "已完成",
}

func (p OrderRepairStatus) ToMap() any {
	return OrderRepairStatusMap
}

func (p OrderRepairStatus) InMap() error {
	if _, ok := OrderRepairStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
