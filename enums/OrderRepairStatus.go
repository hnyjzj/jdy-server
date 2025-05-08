package enums

import (
	"errors"
	"slices"
)

/* 维修单状态 */
// 待付款、已取消、门店已收货、已退款、已送出维修、维修中、已维修送回、待取货、已完成
type OrderRepairStatus int

const (
	OrderRepairStatusWaitPay       OrderRepairStatus = iota + 1 // 待付款
	OrderRepairStatusCancel                                     // 已取消
	OrderRepairStatusStoreReceived                              // 门店已收货
	OrderRepairStatusRefund                                     // 已退款
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
	OrderRepairStatusRefund:        "已退款",
	OrderRepairStatusSendOut:       "已送出维修",
	OrderRepairStatusRepairing:     "维修中",
	OrderRepairStatusSendBack:      "已维修送回",
	OrderRepairStatusWaitPickUp:    "待取货",
	OrderRepairStatusComplete:      "已完成",
}

func (p OrderRepairStatus) CanOperationTo(status OrderRepairStatus) bool {
	transitions := map[OrderRepairStatus][]OrderRepairStatus{
		// 待付款
		OrderRepairStatusWaitPay: {
			OrderRepairStatusCancel,        // 已取消
			OrderRepairStatusStoreReceived, // 门店已收货
		},
		// 已取消
		OrderRepairStatusCancel: {
			OrderRepairStatusWaitPay, // 待付款
		},
		// 门店已收货
		OrderRepairStatusStoreReceived: {
			OrderRepairStatusSendOut, // 已送出维修
			OrderRepairStatusRefund,  // 已退款
		},
		// 已送出维修
		OrderRepairStatusSendOut: {
			OrderRepairStatusRepairing, // 维修中
		},
		// 维修中
		OrderRepairStatusRepairing: {
			OrderRepairStatusSendBack, // 已维修送回
		},
		// 已维修送回
		OrderRepairStatusSendBack: {
			OrderRepairStatusWaitPickUp, // 待取货
		},
		// 待取货
		OrderRepairStatusWaitPickUp: {
			OrderRepairStatusComplete, // 已完成
		},
	}

	if allowed, ok := transitions[p]; ok {
		if slices.Contains(allowed, status) {
			return true
		}
	}
	return true
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
