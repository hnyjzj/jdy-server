package enums

import (
	"errors"
	"slices"
)

/* 产品状态 */
// 正常、已报损、已调拨、已出售、已定出、盘点中、无库存
type ProductStatus int

const (
	ProductStatusNormal   ProductStatus = iota + 1 // 正常
	ProductStatusDamage                            // 已报损
	ProductStatusAllocate                          // 已调拨
	ProductStatusSold                              // 已出售
	ProductStatusReturn                            // 已定出
	ProductStatusCheck                             // 盘点中
	ProductStatusNoStock                           // 无库存
)

var ProductStatusMap = map[ProductStatus]string{
	ProductStatusNormal:   "正常",
	ProductStatusDamage:   "已报损",
	ProductStatusAllocate: "已调拨",
	ProductStatusSold:     "已出售",
	ProductStatusReturn:   "已定出",
	ProductStatusCheck:    "盘点中",
	ProductStatusNoStock:  "无库存",
}

// 判断状态是否可以转换
func (p ProductStatus) CanTransitionTo(newStatus ProductStatus) error {
	transitions := map[ProductStatus][]ProductStatus{
		// 正常
		ProductStatusNormal: {
			ProductStatusDamage,   // 已报损
			ProductStatusAllocate, // 已调拨
			ProductStatusSold,     // 已出售
			ProductStatusReturn,   // 已定出
			ProductStatusCheck,    // 盘点中
			ProductStatusNoStock,  // 无库存
		},
		// 已报损
		ProductStatusDamage: {
			ProductStatusNormal,   // 正常
			ProductStatusAllocate, // 已调拨
		},
		// 已调拨
		ProductStatusAllocate: {
			ProductStatusNormal, // 正常
			ProductStatusDamage, // 已报损
		},
		// 已出售
		ProductStatusSold: {
			ProductStatusNormal, // 正常
			ProductStatusReturn, // 已定出
		},
		// 已定出
		ProductStatusReturn: {
			ProductStatusNormal, // 正常
			ProductStatusDamage, // 已报损
		},
		// 盘点中
		ProductStatusCheck: {
			ProductStatusNormal, // 正常
		},
	}
	if allowed, ok := transitions[p]; ok {
		if slices.Contains(allowed, newStatus) {
			return nil
		}
	}
	return errors.New("非法的状态转换")
}

func (p ProductStatus) ToMap() any {
	return ProductStatusMap
}

func (p ProductStatus) InMap() error {
	if _, ok := ProductStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
