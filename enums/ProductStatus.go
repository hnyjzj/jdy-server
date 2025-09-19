package enums

import (
	"errors"
	"slices"
)

/* 产品状态 */
// 草稿、在库、已报损、调拨中、已出售、已定出、盘点中、无库存
type ProductStatus int

const (
	ProductStatusDraft    ProductStatus = iota + 1 // 草稿
	ProductStatusNormal                            // 在库
	ProductStatusDamage                            // 已报损
	ProductStatusAllocate                          // 调拨中
	ProductStatusSold                              // 已出售
	ProductStatusReturn                            // 已定出
	ProductStatusCheck                             // 盘点中
	ProductStatusNoStock                           // 无库存
)

var ProductStatusMap = map[ProductStatus]string{
	ProductStatusDraft:    "草稿",
	ProductStatusNormal:   "在库",
	ProductStatusDamage:   "已报损",
	ProductStatusAllocate: "调拨中",
	ProductStatusSold:     "已出售",
	ProductStatusReturn:   "已定出",
	ProductStatusCheck:    "盘点中",
	ProductStatusNoStock:  "无库存",
}

// 判断状态是否可以转换
func (p ProductStatus) CanTransitionTo(newStatus ProductStatus) error {
	transitions := map[ProductStatus][]ProductStatus{
		// 在库
		ProductStatusNormal: {
			ProductStatusDamage,   // 已报损
			ProductStatusAllocate, // 调拨中
			ProductStatusSold,     // 已出售
			ProductStatusReturn,   // 已定出
			ProductStatusCheck,    // 盘点中
			ProductStatusNoStock,  // 无库存
		},
		// 已报损
		ProductStatusDamage: {
			ProductStatusNormal,   // 在库
			ProductStatusAllocate, // 调拨中
		},
		// 调拨中
		ProductStatusAllocate: {
			ProductStatusNormal, // 在库
			ProductStatusDamage, // 已报损
		},
		// 已出售
		ProductStatusSold: {
			ProductStatusNormal, // 在库
			ProductStatusReturn, // 已定出
		},
		// 已定出
		ProductStatusReturn: {
			ProductStatusNormal, // 在库
			ProductStatusDamage, // 已报损
		},
		// 盘点中
		ProductStatusCheck: {
			ProductStatusNormal, // 在库
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
