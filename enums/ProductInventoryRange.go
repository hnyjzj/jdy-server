package enums

import (
	"errors"
)

/* 盘点范围 */
// 大类、按材质类型
type ProductInventoryRange int

const (
	ProductInventoryRangeBigType      ProductInventoryRange = iota + 1 // 大类
	ProductInventoryRangeMaterialType                                  // 按材质类型
)

var ProductInventoryRangeMap = map[ProductInventoryRange]string{
	ProductInventoryRangeBigType:      "大类",
	ProductInventoryRangeMaterialType: "按材质类型",
}

func (p ProductInventoryRange) ToMap() any {
	return ProductInventoryRangeMap
}

func (p ProductInventoryRange) InMap() error {
	if _, ok := ProductInventoryRangeMap[p]; !ok {
		return errors.New("not in enum")
	}
	return nil
}
