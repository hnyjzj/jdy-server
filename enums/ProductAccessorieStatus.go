package enums

import (
	"errors"
	"slices"
)

/* 配件状态 */
// 草稿、正常、已调拨、无库存
type ProductAccessorieStatus int

const (
	ProductAccessorieStatusDraft    ProductAccessorieStatus = iota + 1 // 草稿
	ProductAccessorieStatusNormal                                      // 正常
	ProductAccessorieStatusAllocate                                    // 已调拨
	ProductAccessorieStatusNoStock                                     // 无库存
)

var ProductAccessorieStatusMap = map[ProductAccessorieStatus]string{
	ProductAccessorieStatusDraft:    "草稿",
	ProductAccessorieStatusNormal:   "正常",
	ProductAccessorieStatusAllocate: "已调拨",
	ProductAccessorieStatusNoStock:  "无库存",
}

// 判断状态是否可以转换
func (p ProductAccessorieStatus) CanTransitionTo(newStatus ProductAccessorieStatus) error {
	transitions := map[ProductAccessorieStatus][]ProductAccessorieStatus{
		// 正常
		ProductAccessorieStatusNormal: {
			ProductAccessorieStatusAllocate, // 已调拨
			ProductAccessorieStatusNoStock,  // 无库存
		},
		// 已调拨
		ProductAccessorieStatusAllocate: {
			ProductAccessorieStatusNormal, // 正常
		},
	}
	if allowed, ok := transitions[p]; ok {
		if slices.Contains(allowed, newStatus) {
			return nil
		}
	}
	return errors.New("非法的状态转换")
}

func (p ProductAccessorieStatus) ToMap() any {
	return ProductAccessorieStatusMap
}

func (p ProductAccessorieStatus) InMap() error {
	if _, ok := ProductAccessorieStatusMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
