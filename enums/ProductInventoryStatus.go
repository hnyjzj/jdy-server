package enums

import (
	"errors"
)

/* 盘点状态 */
// 草稿、盘点中、待验证、盘点完成、盘点异常、盘点取消
type ProductInventoryStatus int

const (
	ProductInventoryStatusDraft        ProductInventoryStatus = iota + 1 // 草稿
	ProductInventoryStatusInventorying                                   // 盘点中
	ProductInventoryStatusToBeVerified                                   // 待验证
	ProductInventoryStatusCompleted                                      // 盘点完成
	ProductInventoryStatusAbnormal                                       // 盘点异常
	ProductInventoryStatusCancelled                                      // 盘点取消
)

var ProductInventoryStatusMap = map[ProductInventoryStatus]string{
	ProductInventoryStatusDraft:        "草稿",
	ProductInventoryStatusInventorying: "盘点中",
	ProductInventoryStatusToBeVerified: "待验证",
	ProductInventoryStatusCompleted:    "盘点完成",
	ProductInventoryStatusAbnormal:     "盘点异常",
	ProductInventoryStatusCancelled:    "盘点取消",
}

func (p ProductInventoryStatus) ToMap() any {
	return ProductInventoryStatusMap
}

func (p ProductInventoryStatus) InMap() error {
	if _, ok := ProductInventoryStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}

func (p ProductInventoryStatus) String() string {
	return ProductInventoryStatusMap[p]
}

// 判断状态是否可以转换
func (p ProductInventoryStatus) CanTransitionTo(n ProductInventoryStatus) error {
	transitions := map[ProductInventoryStatus][]ProductInventoryStatus{
		ProductInventoryStatusDraft: { // 草稿->
			ProductInventoryStatusInventorying, // 盘点中
			ProductInventoryStatusCancelled,    // 盘点取消
		},
		ProductInventoryStatusInventorying: { // 盘点中->
			ProductInventoryStatusToBeVerified, // 待验证
			ProductInventoryStatusAbnormal,     // 盘点异常
			ProductInventoryStatusCancelled,    // 盘点取消
		},
		ProductInventoryStatusToBeVerified: { // 待验证->
			ProductInventoryStatusCompleted, // 盘点完成
			ProductInventoryStatusAbnormal,  // 盘点异常
			ProductInventoryStatusCancelled, // 盘点取消
		},
		ProductInventoryStatusAbnormal: { // 盘点异常->
			ProductInventoryStatusInventorying, // 盘点中
			ProductInventoryStatusCancelled,    // 盘点取消
		},
		ProductInventoryStatusCancelled: { // 盘点取消->
			ProductInventoryStatusInventorying, // 盘点中
		},
	}

	if allowed, ok := transitions[p]; ok {
		for _, o := range allowed {
			if o == n {
				return nil
			}
		}
	}

	return errors.New("非法的状态转换")
}
